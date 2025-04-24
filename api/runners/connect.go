package runners

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	chserver "github.com/davidebianchi03/chisel/server"
	chsettings "github.com/davidebianchi03/chisel/share/settings"
	dbconn "github.com/davidebianchi03/codebox/db/connection"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

var lock = &sync.Mutex{}

func HandleRunnerConnect(c *gin.Context) {
	runnerId, _ := c.Params.Get("runnerId")

	requestToken := c.Request.Header.Get("X-Codebox-Runner-Token")

	var runner models.Runner
	if err := dbconn.DB.Find(
		&runner,
		map[string]interface{}{
			"id": runnerId,
		},
	).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if runner.Token != requestToken {
		c.JSON(http.StatusUnauthorized, gin.H{
			"detail": "missing or invalid token",
		})
		return
	}

	if c.Request.Method == http.MethodPost {
		// assign a free port on the local host to the runner
		minPort := 20000
		maxPort := 50000

		assignedPort := 0
		lock.Lock()

		for i := minPort; i < maxPort; i++ {
			var count int64
			r := dbconn.DB.Model(&models.Runner{}).Where(map[string]interface{}{
				"port": i,
			}).Count(&count)

			if r.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"detail": "internal server error",
				})
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

			c.JSON(http.StatusOK, gin.H{
				"port": assignedPort,
			})
		} else {
			log.Println("there are no free ports available")
			c.JSON(http.StatusTeapot, gin.H{
				"detail": "no free ports available",
			})
		}
		lock.Unlock()
	} else {
		// forward port using chisel
		serverConfig := chserver.Config{
			Reverse: true,
			AuthCallback: func(r *chsettings.Remote) bool {
				return r.LocalPort == strconv.Itoa(int(runner.Port))
			},
		}
		s, err := chserver.NewServer(&serverConfig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"detail": "internal server error",
			})
			return
		}
		s.Debug = false
		s.HandleClientHandler(c.Writer, c.Request)

		// release the port
		runner.Port = 0
		dbconn.DB.Save(&runner)
	}
}
