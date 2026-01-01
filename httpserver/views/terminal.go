package views

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
)

func HandleTerminalView(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "errors.html", gin.H{
			"title":   "Unknown error",
			"message": "Internal server error",
		})
		return
	}

	workspaceId, err := utils.GetUIntParamFromContext(ctx, "workspaceId")
	if err != nil {
		ctx.HTML(http.StatusNotFound, "errors.html", gin.H{
			"title":   "Not found",
			"message": "Workspace not found",
		})
		return
	}

	containerName, found := ctx.Params.Get("containerName")
	if !found {
		ctx.HTML(http.StatusNotFound, "errors.html", gin.H{
			"title":   "Not found",
			"message": "Container not found",
		})
		return
	}

	workspace, err := models.RetrieveWorkspaceByUserAndId(user, workspaceId)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "errors.html", gin.H{
			"title":   "Unknown error",
			"message": "Internal server error",
		})
		return
	}

	if workspace == nil {
		ctx.HTML(http.StatusNotFound, "errors.html", gin.H{
			"title":   "Not found",
			"message": "Workspace not found",
		})
		return
	}

	container, err := models.RetrieveWorkspaceContainerByName(*workspace, containerName)
	if err != nil {
		ctx.HTML(http.StatusNotFound, "errors.html", gin.H{
			"title":   "Unknown error",
			"message": "Internal server error",
		})
		return
	}

	if container == nil {
		ctx.HTML(http.StatusNotFound, "errors.html", gin.H{
			"title":   "Not found",
			"message": "Container not found",
		})
		return
	}

	ctx.HTML(http.StatusOK, "terminal.html", gin.H{
		"page_title": fmt.Sprintf(
			"Terminal %s - %s",
			workspace.Name,
			container.ContainerName,
		),
		"terminal_api": fmt.Sprintf(
			"/api/v1/workspace/%d/container/%s/terminal",
			workspace.ID,
			container.ContainerName,
		),
	})
}
