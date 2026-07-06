package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
	"gitlab.com/codebox4073715/codebox/logging"
)

// HandleAdminListSystemLogs godoc
// @Summary List System Logs
// @Schemes
// @Description List all system logs
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} []serializers.SystemLogRowSerializer
// @Router /api/v1/admin/system-logs [get]
func HandleAdminListSystemLogs(c *gin.Context) {
	logs, err := logging.ListAllLogs()
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	c.JSON(200, serializers.LoadMultipleSystemLogRows(logs))
}
