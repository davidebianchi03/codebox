package workspaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
)

/*
Endpoint that handles ssh connection
*/
func HandleForwardSsh(c *gin.Context) {
	container, err := retrieveContainerByWorkspaceAndName(c)
	if err != nil {
		return
	}

	workspace := container.Workspace

	runner, err := models.RetrieveRunnerByID(*workspace.RunnerID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if runner == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "runner not found")
		return
	}

	ri := runnerinterface.RunnerInterface{
		Runner: runner,
	}
	if err := ri.ForwardSsh(&workspace, container, c.Writer, c.Request); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	// If the connection is closed, we return a 418 status code
	utils.ErrorResponse(c, http.StatusTeapot, "connection closed")
}
