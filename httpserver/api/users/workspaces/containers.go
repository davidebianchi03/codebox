package workspaces

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
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

/*
RetrieveWorkspaceContainerFromContext is a helper function that retrieves
a workspace container from the context.
*/
func retrieveWorkspaceContainerFromContext(c *gin.Context) (*models.WorkspaceContainer, error) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		utils.ErrorResponse(
			c, http.StatusInternalServerError, "internal server error",
		)
		return nil, err
	}

	workspaceId, err := utils.GetUIntParamFromContext(c, "workspaceId")
	if err != nil {
		utils.ErrorResponse(
			c, http.StatusNotFound, "workspace not found",
		)
		return nil, errors.New("workspace id not found in path")
	}

	containerName, found := c.Params.Get("containerName")
	if !found {
		utils.ErrorResponse(
			c, http.StatusNotFound, "container not found",
		)
		return nil, errors.New("container name not found in path")
	}

	workspace, err := models.RetrieveWorkspaceByUserAndId(user, workspaceId)
	if err != nil {
		utils.ErrorResponse(
			c, http.StatusInternalServerError, "internal server error",
		)
		return nil, err
	}

	if workspace == nil {
		utils.ErrorResponse(
			c, http.StatusNotFound, "workspace not found",
		)
		return nil, errors.New("workspace not found")
	}

	container, err := models.RetrieveWorkspaceContainerByName(*workspace, containerName)
	if err != nil {
		utils.ErrorResponse(
			c, http.StatusInternalServerError, "internal server error",
		)
		return nil, err
	}

	if container == nil {
		utils.ErrorResponse(
			c, http.StatusNotFound, "container not found",
		)
		return nil, errors.New("container not found")
	}

	return container, nil
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
func RetrieveWorkspaceContainersByWorkspace(c *gin.Context) {
	container, err := retrieveWorkspaceContainerFromContext(c)
	if err != nil {
		return
	}

	c.JSON(
		http.StatusOK,
		serializers.LoadWorkspaceContainerSerializer(container),
	)
}

// WorkspaceContainerListDirectory godoc
// @Summary WorkspaceContainerListDirectory
// @Schemes
// @Description List the contents of a directory in a container.
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param path query string true "Directory path"
// @Success 200 {object} []serializers.ContainerFileInfoSerializer
// @Failure 400 {object} serializers.ErrorSerializer "Bad request (e.g., missing or invalid 'path' parameter, or provided path is not a directory)"
// @Failure 403 {object} serializers.ErrorSerializer "Forbidden (e.g., permission denied when trying to access the directory)"
// @Failure 404 {object} serializers.ErrorSerializer "workspace, container or requested path not found"
// @Failure 500 {object} serializers.ErrorSerializer "Internal server error"
// @Router /api/v1/workspace/:workspaceId/container/:containerName/fs/list-directory [get]
func WorkspaceContainerListDirectory(c *gin.Context) {
	container, err := retrieveWorkspaceContainerFromContext(c)
	if err != nil {
		return
	}

	ri := runnerinterface.RunnerInterface{
		Runner: container.Workspace.Runner,
	}

	path := c.Query("path")
	if path == "" {
		utils.ErrorResponse(
			c, http.StatusBadRequest, "path query parameter is required",
		)
		return
	}

	files, err := ri.ContainerFsListDir(
		&container.Workspace,
		container,
		path,
	)
	if err != nil {
		if runnerinterface.IsPathNotExist(err) {
			utils.ErrorResponse(
				c, http.StatusNotFound, err.Error(),
			)
		} else if runnerinterface.IsPermissionDenied(err) {
			utils.ErrorResponse(
				c, http.StatusForbidden, err.Error(),
			)
		} else if runnerinterface.IsPathIsNotADir(err) {
			utils.ErrorResponse(
				c, http.StatusBadRequest, err.Error(),
			)
		} else {
			// TODO: log error
			utils.ErrorResponse(
				c, http.StatusInternalServerError, "internal server error",
			)
		}
		return
	}

	c.JSON(
		http.StatusOK,
		serializers.LoadMultipleContainerFileInfoSerializers(files),
	)
}

// WorkspaceContainerGetItemInfo godoc
// @Summary WorkspaceContainerGetItemInfo
// @Schemes
// @Description Get detailed information about a file or directory
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param path query string true "File or directory path"
// @Success 200 {object} serializers.ContainerFileInfoSerializer
// @Failure 400 {object} serializers.ErrorSerializer "Bad request (e.g., missing or invalid 'path' parameter)"
// @Failure 403 {object} serializers.ErrorSerializer "Forbidden (permission denied)"
// @Failure 404 {object} serializers.ErrorSerializer "workspace, container or requested path not found"
// @Failure 500 {object} serializers.ErrorSerializer "Internal server error"
// @Router /api/v1/workspace/:workspaceId/container/:containerName/fs/get-item-info [get]
func WorkspaceContainerGetItemInfo(c *gin.Context) {
	container, err := retrieveWorkspaceContainerFromContext(c)
	if err != nil {
		return
	}

	ri := runnerinterface.RunnerInterface{
		Runner: container.Workspace.Runner,
	}

	path := c.Query("path")
	if path == "" {
		utils.ErrorResponse(
			c, http.StatusBadRequest, "path query parameter is required",
		)
		return
	}

	files, err := ri.ContainerFsGetItemInfo(
		&container.Workspace,
		container,
		path,
	)
	if err != nil {
		if runnerinterface.IsPathNotExist(err) {
			utils.ErrorResponse(
				c, http.StatusNotFound, err.Error(),
			)
		} else if runnerinterface.IsPermissionDenied(err) {
			utils.ErrorResponse(
				c, http.StatusForbidden, err.Error(),
			)
		} else if runnerinterface.IsPathIsNotADir(err) {
			utils.ErrorResponse(
				c, http.StatusBadRequest, err.Error(),
			)
		} else {
			// TODO: log error
			utils.ErrorResponse(
				c, http.StatusInternalServerError, "internal server error",
			)
		}
		return
	}

	c.JSON(
		http.StatusOK,
		serializers.LoadContainerFileInfoSerializer(files),
	)
}

// WorkspaceContainerGetItemInfo godoc
// @Summary WorkspaceContainerGetItemInfo
// @Schemes
// @Description Get detailed information about a file or directory
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param path query string true "File or directory path"
// @Success 200 {object} serializers.ContainerFileInfoSerializer
// @Failure 400 {object} serializers.ErrorSerializer "Bad request (e.g., missing or invalid 'path' parameter)"
// @Failure 403 {object} serializers.ErrorSerializer "Forbidden (permission denied)"
// @Failure 404 {object} serializers.ErrorSerializer "workspace, container or requested path not found"
// @Failure 500 {object} serializers.ErrorSerializer "Internal server error"
// @Router /api/v1/workspace/:workspaceId/container/:containerName/fs/get-item-info [get]
func WorkspaceContainerCreateDirectory(c *gin.Context) {
	container, err := retrieveWorkspaceContainerFromContext(c)
	if err != nil {
		return
	}

	ri := runnerinterface.RunnerInterface{
		Runner: container.Workspace.Runner,
	}

	path := c.Query("path")
	if path == "" {
		utils.ErrorResponse(
			c, http.StatusBadRequest, "path query parameter is required",
		)
		return
	}

	files, err := ri.ContainerFsGetItemInfo(
		&container.Workspace,
		container,
		path,
	)
	if err != nil {
		if runnerinterface.IsPathNotExist(err) {
			utils.ErrorResponse(
				c, http.StatusNotFound, err.Error(),
			)
		} else if runnerinterface.IsPermissionDenied(err) {
			utils.ErrorResponse(
				c, http.StatusForbidden, err.Error(),
			)
		} else if runnerinterface.IsPathIsNotADir(err) {
			utils.ErrorResponse(
				c, http.StatusBadRequest, err.Error(),
			)
		} else {
			// TODO: log error
			utils.ErrorResponse(
				c, http.StatusInternalServerError, "internal server error",
			)
		}
		return
	}

	c.JSON(
		http.StatusOK,
		serializers.LoadContainerFileInfoSerializer(files),
	)
}
