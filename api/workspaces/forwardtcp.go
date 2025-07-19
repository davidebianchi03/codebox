package workspaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
)

func HandleForwardTcp(ctx *gin.Context) {
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

	portNumber, found := ctx.Params.Get("portNumber")
	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	var workspace models.Workspace
	result := dbconn.DB.
		Preload("Runner").
		Where(
			map[string]interface{}{
				"ID": workspaceId,
			}).
		Find(&workspace)

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

	if workspace.Runner == nil {
		ctx.JSON(http.StatusTeapot, gin.H{
			"detail": "runner not found",
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

	// retrieve port
	port := models.WorkspaceContainerPort{}
	r = dbconn.DB.Model(&models.WorkspaceContainerPort{}).First(&port, map[string]interface{}{
		"port_number":  portNumber,
		"container_id": container.ID,
	})

	if r.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if r.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "port is not exposed",
		})
		return
	}

	ri := runnerinterface.RunnerInterface{
		Runner: workspace.Runner,
	}
	if err := ri.ForwardTcpPort(&workspace, &container, ctx.Writer, ctx.Request, "1234"); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusTeapot, gin.H{
		"detail": "connection closed",
	})
}
