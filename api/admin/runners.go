package admin

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gocraft/work"
	"gitlab.com/codebox4073715/codebox/api/serializers"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/bgtasks"
	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
)

// AdminRunners godoc
// @Summary List all available runners
// @Schemes
// @Description List all available runners
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} serializers.AdminRunnersSerializer[]
// @Router /api/v1/admin/runners [get]
func HandleAdminListRunners(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "-1"
	}

	// validate limit
	parsedLimit, err := strconv.Atoi(limit)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid limit")
		return
	}

	if parsedLimit < -1 || parsedLimit == 0 {
		utils.ErrorResponse(c, 400, "invalid limit")
		return
	}

	runners, err := models.ListRunners(parsedLimit, 0)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, serializers.LoadMultipleAdminRunnerSerializer(runners))
}

// HandleAdminRetrieveRunners godoc
// @Summary Retrive a runner by its id
// @Schemes
// @Description Retrive a runner by its id
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} serializers.AdminRunnersSerializer
// @Router /api/v1/admin/runners/:id [get]
func HandleAdminRetrieveRunners(c *gin.Context) {
	runnerId, _ := c.Params.Get("runnerId")

	id, err := strconv.Atoi(runnerId)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "runner not found")
		return
	}

	runner, err := models.RetrieveRunnerByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if runner == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "runner not found")
		return
	}

	c.JSON(http.StatusOK, serializers.LoadAdminRunnerSerializer(runner))
}

type HandleAdminCreateRunnerRequestBody struct {
	Name         string `json:"name" binding:"required"`
	Type         string `json:"type" binding:"required"`
	UsePublicUrl bool   `json:"use_public_url"`
	PublicUrl    string `json:"public_url"`
}

// HandleAdminCreateRunner godoc
// @Summary Create a runner
// @Schemes
// @Description Create a runner
// @Tags Admin
// @Accept json
// @Produce json
// @Success 201
// @Param request body HandleAdminCreateRunnerRequestBody true "Runner details"
// @Router /api/v1/admin/runners [post]
func HandleAdminCreateRunner(c *gin.Context) {
	// parse and validate request body
	var parsedBody HandleAdminCreateRunnerRequestBody
	err := c.ShouldBindBodyWithJSON(&parsedBody)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	runner, err := models.RetrieveRunnerByName(parsedBody.Name)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if runner != nil {
		utils.ErrorResponse(c, http.StatusConflict, "another runner with the same name already exists")
		return
	}

	runnerTypeFound := false
	for _, rt := range config.ListAvailableRunnerTypes() {
		if rt.ID == parsedBody.Type {
			runnerTypeFound = true
		}
	}

	if !runnerTypeFound {
		utils.ErrorResponse(c, http.StatusConflict, "runner type not found")
		return
	}

	if parsedBody.UsePublicUrl {
		if parsedBody.PublicUrl == "" {
			utils.ErrorResponse(c, http.StatusConflict, "'public_url' is required")
			return
		}

		exists, err := models.DoesRunnerExistWithUrl(parsedBody.PublicUrl)
		if err != nil {
			log.Println(err)
			utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
			return
		}

		if exists {
			utils.ErrorResponse(
				c,
				http.StatusConflict,
				"another runner with the same public url already exists",
			)
			return
		}
	}

	// create the runner
	runner, err = models.CreateRunner(
		parsedBody.Name,
		parsedBody.Type,
		parsedBody.UsePublicUrl,
		parsedBody.PublicUrl,
	)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	bgtasks.BgTasksEnqueuer.Enqueue("ping_runners", work.Q{})

	c.JSON(http.StatusCreated, gin.H{
		"id":    runner.ID,
		"token": runner.Token,
	})
}

type AdminUpdateRunnerRequestBody struct {
	Name         string `json:"name" binding:"required"`
	Type         string `json:"type" binding:"required"`
	UsePublicUrl *bool  `json:"use_public_url" binding:"required"`
	PublicUrl    string `json:"public_url" binding:"required"`
}

// HandleAdminCreateRunner godoc
// @Summary Create a runner
// @Schemes
// @Description Create a runner
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} serializers.AdminRunnersSerializer
// @Param request body AdminUpdateRunnerRequestBody true "Runner details"
// @Router /api/v1/admin/runners/:id [put]
func HandleAdminUpdateRunner(c *gin.Context) {
	runnerId, _ := c.Params.Get("runnerId")

	id, err := strconv.Atoi(runnerId)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "runner not found")
	}

	runner, err := models.RetrieveRunnerByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if runner == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "runner not found")
		return
	}

	var reqBody AdminUpdateRunnerRequestBody
	err = c.ShouldBindBodyWithJSON(&reqBody)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing or invalid field")
		return
	}

	runner.Name = reqBody.Name
	runner.Type = reqBody.Type
	runner.UsePublicUrl = *reqBody.UsePublicUrl
	runner.PublicUrl = reqBody.PublicUrl

	if err := models.UpdateRunner(*runner); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, serializers.LoadAdminRunnerSerializer(runner))
}
