package runners

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	chserver "github.com/davidebianchi03/chisel/server"
	chsettings "github.com/davidebianchi03/chisel/share/settings"
	"github.com/gin-gonic/gin"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
	"gitlab.com/codebox4073715/codebox/logging"
)

var lock = &sync.Mutex{}

// RunnerRequestPort godoc
// @Summary API used by runners to request a free port to use on server
// @Schemes
// @Description API used by runners to request a free port to use on server
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200
func HandleRunnerRequestPort(c *gin.Context) {
	runnerId, _ := utils.GetUIntParamFromContext(c, "runnerId")
	runner, err := models.RetrieveRunnerByID(runnerId)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	if runner == nil {
		utils.ErrorResponse(
			c,
			http.StatusNotFound,
			"runner not found",
		)
		return
	}

	// assign a free port on the local host to the runner
	minPort := 20000
	maxPort := 50000

	assignedPort := 0
	lock.Lock()
	defer lock.Unlock()

	for i := minPort; i < maxPort; i++ {
		var count int64
		r := dbconn.DB.Model(&models.Runner{}).Where(map[string]interface{}{
			"port": i,
		}).Count(&count)

		if r.Error != nil {
			utils.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"internal server error",
			)
			return
		}

		if count == 0 && int(runner.Port) != i {
			assignedPort = i
			break
		}
	}

	if assignedPort != 0 {
		// append agent to the list of connected agents
		runner.Port = uint(assignedPort)
		dbconn.DB.Save(&runner)

		logging.Info(
			"port %d has been assigned to runner %d",
			runner.Port,
			runner.ID,
		)

		c.JSON(http.StatusOK, gin.H{
			"port": assignedPort,
		})
	} else {
		logging.Error(
			"cannot assign a port to runner %d, no free ports available",
			runner.ID,
		)

		utils.ErrorResponse(
			c,
			http.StatusTeapot,
			"no free ports available",
		)
	}
}

// RunenrConnect godoc
// @Summary API used by runners to forward port
// @Schemes
// @Description API used by runners to forward port
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200
func HandleRunnerConnect(c *gin.Context) {
	runnerId, _ := utils.GetUIntParamFromContext(c, "runnerId")
	runner, err := models.RetrieveRunnerByID(runnerId)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	if runner == nil {
		utils.ErrorResponse(
			c,
			http.StatusNotFound,
			"runner not found",
		)
		return
	}

	// forward port using chisel
	serverConfig := chserver.Config{
		Reverse: true,
		AuthCallback: func(r *chsettings.Remote) bool {
			return r.LocalPort == strconv.Itoa(int(runner.Port))
		},
		KeepAlive: time.Second,
	}
	s, err := chserver.NewServer(&serverConfig)
	if err != nil {
		logging.Error(
			"failed to setup port forwarding for runner %d, %s",
			runner.ID,
			err.Error(),
		)

		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	logging.Info(
		"runner %d is connected",
		runner.ID,
	)

	s.Debug = false
	s.HandleClientHandler(c.Writer, c.Request)

	// release the port
	runner, err = models.RetrieveRunnerByID(runnerId)
	if err == nil && runner != nil {
		runner.Port = 0
		dbconn.DB.Save(&runner)
	}

	logging.Warn(
		"runner %d disconnected",
		runner.ID,
	)
}
