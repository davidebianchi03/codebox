package workspaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
)

func HandleForwardSsh(ctx *gin.Context) {
	workspaceId, found := ctx.Params.Get("workspaceId")
	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	containerName, found := ctx.Params.Get("containerName")
	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	var workspace models.Workspace
	result := dbconn.DB.Where(map[string]interface{}{"ID": workspaceId}).Preload("Runner").Find(&workspace)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	// retrieve development container
	container := models.WorkspaceContainer{}
	r := dbconn.DB.First(&container, map[string]interface{}{
		"workspace_id":   workspace.ID,
		"container_name": containerName,
	})

	if r.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if container.ID <= 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "container not found, check that workspace is running and that you can connect to this container",
		})
		return
	}

	if workspace.Runner == nil {
		ctx.JSON(http.StatusTeapot, gin.H{
			"detail": "runner not found",
		})
		return
	}

	ri := runnerinterface.RunnerInterface{
		Runner: workspace.Runner,
	}
	if err := ri.ForwardSsh(&workspace, &container, ctx.Writer, ctx.Request); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusTeapot, gin.H{
		"detail": "connection closed",
	})
}
