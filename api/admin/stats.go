package admin

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/serializers"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/db/models"
)

// AdminStats godoc
// @Summary Admin Stats
// @Schemes
// @Description Admin Stats
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} serializers.AdminStatsSerializer
// @Router /api/v1/admin/stats [get]
func HandleAdminStats(c *gin.Context) {
	loginsCount, err := models.GetLoginCountPerDayInLast7Days()
	if err != nil {
		utils.ErrorResponse(c, 500, "cannot retrieve login stats")
		return
	}

	onlineRunnersCount, err := models.CountOnlineRunners()
	if err != nil {
		utils.ErrorResponse(c, 500, "cannot retrieve online runners count")
		return
	}

	usersCount, err := models.CountAllUsers()
	if err != nil {
		utils.ErrorResponse(c, 500, "cannot retrieve users count")
		return
	}

	onlineWorkspaces, err := models.CountAllOnlineWorkspaces()
	if err != nil {
		utils.ErrorResponse(c, 500, "cannot retrieve online workspaces count")
		return
	}

	c.JSON(200, serializers.LoadAdminStatsSerializer(
		loginsCount,
		usersCount,
		onlineRunnersCount,
		onlineWorkspaces,
	))
}
