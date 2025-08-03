package workspaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
)

func HandleForwardSsh(ctx *gin.Context) {
	container, err := retrieveContainerByWorkspaceAndName(ctx)
	if err != nil {
		return
	}

	workspace := container.Workspace

	if workspace.Runner == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "runner not found")
		return
	}

	ri := runnerinterface.RunnerInterface{
		Runner: workspace.Runner,
	}
	if err := ri.ForwardSsh(&workspace, container, ctx.Writer, ctx.Request); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	// If the connection is closed, we return a 418 status code
	utils.ErrorResponse(ctx, http.StatusTeapot, "connection closed")
}
