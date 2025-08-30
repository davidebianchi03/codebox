package workspaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
)

func HandleForwardTcp(ctx *gin.Context) {
	container, err := retrieveContainerByWorkspaceAndName(ctx)
	if err != nil {
		return
	}

	workspace := container.Workspace

	runner, err := models.RetrieveRunnerByID(workspace.RunnerID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if runner == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "runner not found")
		return
	}

	// retrieve port number from URL parameter
	portNumber, err := utils.GetUIntParamFromContext(ctx, "portNumber")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid port number")
		return
	}

	port, err := models.RetrieveContainerPortByPortNumber(*container, portNumber)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if port == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "port not found",
		})
		return
	}

	ri := runnerinterface.RunnerInterface{
		Runner: runner,
	}
	if err := ri.ForwardTcpPort(&workspace, container, ctx.Writer, ctx.Request, portNumber); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.ErrorResponse(ctx, http.StatusTeapot, "connection closed")
}
