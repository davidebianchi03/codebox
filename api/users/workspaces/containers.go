package workspaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/db/models"
)

// ListWorkspaceContainersByWorkspace godoc
// @Summary ListWorkspaceContainersByWorkspace
// @Schemes
// @Description List all containers for a workspace
// @Tags Workspaces
// @Accept json
// @Produce json
// @Success 200 {object} []serializers.WorkspaceContainerSerializer
// @Router /api/v1/workspace/:workspaceId/container [get]
func ListWorkspaceContainersByWorkspace(ctx *gin.Context) {
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

	containers, err := models.ListWorkspaceContainersByWorkspace(*workspace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	ctx.JSON(
		http.StatusOK,
		serializers.LoadMultipleWorkspaceContainerSerializers(containers),
	)
}

// RetrieveWorkspaceContainersByWorkspace godoc
// @Summary RetrieveWorkspaceContainersByWorkspace
// @Schemes
// @Description Retrieve a specific container by name in a workspace
// @Tags Workspaces
// @Accept json
// @Produce json
// @Success 200 {object} serializers.WorkspaceContainerSerializer
// @Router /api/v1/workspace/:workspaceId/container/:containerName [get]
func RetrieveWorkspaceContainersByWorkspace(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	workspaceId, err := utils.GetUIntParamFromContext(ctx, "workspaceId")
	if err != nil {
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

	workspace, err := models.RetrieveWorkspaceByUserAndId(user, workspaceId)
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

	container, err := models.RetrieveWorkspaceContainerByName(*workspace, containerName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if container == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "container not found",
		})
		return
	}

	ctx.JSON(
		http.StatusOK,
		serializers.LoadWorkspaceContainerSerializer(container),
	)
}
