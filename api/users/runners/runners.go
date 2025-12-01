package runners

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/db/models"
)

// List runners godoc
// @Summary List runners
// @Schemes
// @Description List runners
// @Tags Runners
// @Accept json
// @Produce json
// @Success 200 {object} []serializers.RunnerSerializer
// @Router /api/v1/runners [get]
func HandleListRunners(c *gin.Context) {
	runners, err := models.ListRunners(-1, 0)
	if err != nil {
		utils.ErrorResponse(
			c, http.StatusInternalServerError, "internal server error",
		)
		return
	}

	c.JSON(http.StatusOK, serializers.LoadMultipleRunnerSerializer(runners))
}
