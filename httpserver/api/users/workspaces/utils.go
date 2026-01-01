package workspaces

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
)

// retrieveContainerByWorkspaceAndName retrieves a container by
// workspace ID and container name from the context.
// It returns the container if found, or an error if not
// found or if there is an internal error.
func retrieveContainerByWorkspaceAndName(ctx *gin.Context) (*models.WorkspaceContainer, error) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "internal server error")
		return nil, errors.New("unknown error")
	}

	workspaceId, err := utils.GetUIntParamFromContext(ctx, "workspaceId")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "workspace not found")
		return nil, errors.New("workspace not found")
	}

	containerName, found := ctx.Params.Get("containerName")
	if !found {
		utils.ErrorResponse(ctx, http.StatusNotFound, "container not found")
		return nil, errors.New("container not found")
	}

	workspace, err := models.RetrieveWorkspaceByUserAndId(user, workspaceId)
	if err != nil || workspace == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "workspace not found")
		return nil, errors.New("workspace not found")
	}

	container, err := models.RetrieveWorkspaceContainerByName(*workspace, containerName)
	if err != nil || container == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "container not found")
		return nil, errors.New("container not found")
	}

	return container, nil
}
