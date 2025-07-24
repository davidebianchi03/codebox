package workspaces

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/serializers"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/db/models"
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

// ListContainerPortsByWorkspaceContainer godoc
// @Summary ListContainerPortsByWorkspaceContainer
// @Schemes
// @Description List all ports for a container in a workspace
// @Tags Workspaces
// @Accept json
// @Produce json
// @Success 200 {object} []serializers.WorkspaceContainerPort
// @Router /api/v1/workspace/:workspaceId/container/:containerName/port [get]
func ListContainerPortsByWorkspaceContainer(c *gin.Context) {
	container, err := retrieveContainerByWorkspaceAndName(c)
	if err != nil {
		return
	}

	containerPorts, err := models.ListContainerPortsByWorkspaceContainer(*container)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(
		http.StatusOK,
		serializers.LoadMultipleWorkspaceContainerPorts(containerPorts),
	)
}

// RetrieveContainerPortsByWorkspaceContainer godoc
// @Summary RetrieveContainerPortsByWorkspaceContainer
// @Schemes
// @Description Retrieve a specific port by number for a container in a workspace
// @Tags Workspaces
// @Accept json
// @Produce json
// @Success 200 {object} serializers.WorkspaceContainerPort
// @Router /api/v1/workspace/:workspaceId/container/:containerName/port/:portNumber [get]
func RetrieveContainerPortsByWorkspaceContainer(ctx *gin.Context) {
	container, err := retrieveContainerByWorkspaceAndName(ctx)
	if err != nil {
		return
	}

	portNumber, err := utils.GetUIntParamFromContext(ctx, "portNumber")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "port not found")
		return
	}

	if portNumber < 1 || portNumber > 65535 {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid port number")
		return
	}
	containerPort, err := models.RetrieveContainerPortByPortNumber(*container, portNumber)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
		return
	}

	if containerPort == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "port not found")
		return
	}

	ctx.JSON(
		http.StatusOK,
		serializers.LoadWorkspaceContainerPort(containerPort),
	)
}

type CreateContainerPortRequestBody struct {
	PortNumber  uint   `json:"port_number" binding:"required"`
	ServiceName string `json:"service_name" binding:"required"`
	Public      bool   `json:"public"`
}

// HandleCreateContainerPortByWorkspaceContainer godoc
// @Summary Expose a new port for a container in a workspace
// @Schemes
// @Description Expose a new port for a container in a workspace
// @Tags Workspaces
// @Accept json
// @Produce json
// @Param request body CreateContainerPortRequestBody true "CreateContainerPortRequestBody	"
// @Success 201 {object} serializers.WorkspaceContainerPort
// @Router/api/v1/workspace/:workspaceId/container/:containerName/port [post]
func HandleCreateContainerPortByWorkspaceContainer(c *gin.Context) {
	container, err := retrieveContainerByWorkspaceAndName(c)
	if err != nil {
		return
	}

	var reqBody CreateContainerPortRequestBody
	if err := c.ShouldBindBodyWithJSON(&reqBody); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing or invalid request argument")
		return
	}

	if reqBody.PortNumber < 1 || reqBody.PortNumber > 65535 {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid port number")
		return
	}

	port, err := models.RetrieveContainerPortByPortNumber(*container, reqBody.PortNumber)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if port != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "this port is already exposed")
		return
	}

	port, err = models.RetrieveContainerPortByServiceName(*container, reqBody.ServiceName)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if port != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "this service name is already exposed")
		return
	}

	containerPort, err := models.CreateContainerPort(
		*container,
		reqBody.ServiceName,
		reqBody.PortNumber,
		reqBody.Public,
	)

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(
		http.StatusCreated,
		serializers.LoadWorkspaceContainerPort(containerPort),
	)
}

// HandleDeleteContainerPortByWorkspaceContainer godoc
// @Summary DeleteContainerPortByWorkspaceContainer
// @Schemes
// @Description Delete a specific port by number for a container in a workspace
// @Tags Workspaces
// @Accept json
// @Produce json
// @Success 204
// @Router/api/v1/workspace/:workspaceId/container/:containerName/port/:portNumber [delete]
func HandleDeleteContainerPortByWorkspaceContainer(c *gin.Context) {
	container, err := retrieveContainerByWorkspaceAndName(c)
	if err != nil {
		return
	}

	portNumber, err := utils.GetUIntParamFromContext(c, "portNumber")
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "port not found")
		return
	}

	port, err := models.RetrieveContainerPortByPortNumber(*container, portNumber)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if port == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "port not found")
		return
	}

	if err := models.DeleteContainerPort(port); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"detail": "port has been removed",
	})
}
