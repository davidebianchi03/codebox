package workspaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
)

// HandleRetrieveWorkspaceLogs godoc
// @Summary Retrieve workspace logs
// @Schemes
// @Description Retrieve workspace logs
// @Tags Workspaces
// @Accept json
// @Produce json
// @Success 200
// @Router /api/v1/workspace/:workspaceId/logs [get]
func HandleRetrieveWorkspaceLogs(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	id, err := utils.GetUIntParamFromContext(ctx, "workspaceId")
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	workspace, err := models.RetrieveWorkspaceByUserAndId(user, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if workspace == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	logs, _ := workspace.RetrieveLogs()
	ctx.JSON(http.StatusOK, gin.H{
		"logs": logs,
	})
}
