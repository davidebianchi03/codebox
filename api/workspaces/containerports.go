package workspaces

import (
	"net/http"

	"github.com/davidebianchi03/codebox/api/utils"
	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

func ListContainerPortsByWorkspaceContainer(c *gin.Context) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	workspaceId, found := c.Params.Get("workspaceId")
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	containerName, found := c.Params.Get("containerName")
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	var workspace models.Workspace
	result := db.DB.Find(&workspace, map[string]interface{}{"ID": workspaceId, "user_id": user.ID})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	var container models.WorkspaceContainer
	result = db.DB.Find(&container, map[string]interface{}{"container_name": containerName, "workspace_id": workspace.ID})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "container not found",
		})
		return
	}

	var containerPorts []models.WorkspaceContainerPort
	result = db.DB.Find(&containerPorts, map[string]interface{}{"container_id": container.ID})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, containerPorts)
}

func RetrieveContainerPortsByWorkspaceContainer(c *gin.Context) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	workspaceId, found := c.Params.Get("workspaceId")
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	containerName, found := c.Params.Get("containerName")
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	portNumber, found := c.Params.Get("portNumber")
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	var workspace models.Workspace
	result := db.DB.Find(&workspace, map[string]interface{}{"ID": workspaceId, "user_id": user.ID})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	var container models.WorkspaceContainer
	result = db.DB.Find(&container, map[string]interface{}{"container_name": containerName, "workspace_id": workspace.ID})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "container not found",
		})
		return
	}

	var containerPort models.WorkspaceContainerPort
	result = db.DB.Find(&containerPort, map[string]interface{}{"container_id": container.ID, "port_number": portNumber})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "port not found",
		})
		return
	}

	c.JSON(http.StatusOK, containerPort)
}

func HandleCretateContainerPortByWorkspaceContainer(c *gin.Context) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	workspaceId, found := c.Params.Get("workspaceId")
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	containerName, found := c.Params.Get("containerName")
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	var workspace models.Workspace
	result := db.DB.Find(&workspace, map[string]interface{}{"ID": workspaceId, "user_id": user.ID})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	var container models.WorkspaceContainer
	result = db.DB.Find(&container, map[string]interface{}{"container_name": containerName, "workspace_id": workspace.ID})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "container not found",
		})
		return
	}

	var reqBody struct {
		PortNumber  uint   `json:"port_number" binding:"required"`
		ServiceName string `json:"service_name" binding:"required"`
		Public      bool   `json:"public"`
	}

	if err := c.ShouldBindBodyWithJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing or invalid request argument",
		})
		return
	}

	if reqBody.PortNumber < 1 || reqBody.PortNumber > 65535 {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid port number",
		})
		return
	}

	var count int64

	if err := db.DB.
		Model(&models.WorkspaceContainerPort{}).
		Where(map[string]interface{}{
			"container_id": container.ID,
			"port_number":  reqBody.PortNumber,
		}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "this port is already exposed",
		})
		return
	}

	if err := db.DB.
		Model(&models.WorkspaceContainerPort{}).
		Where(map[string]interface{}{
			"container_id": container.ID,
			"service_name": reqBody.ServiceName,
		}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "another port with the same name already exists",
		})
		return
	}

	containerPort := models.WorkspaceContainerPort{
		ContainerID: container.ID,
		ServiceName: reqBody.ServiceName,
		PortNumber:  uint(reqBody.PortNumber),
		Public:      reqBody.Public,
	}

	db.DB.Save(&containerPort)
	c.JSON(http.StatusCreated, containerPort)
}

func HandleDeleteContainerPortByWorkspaceContainer(c *gin.Context) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	workspaceId, found := c.Params.Get("workspaceId")
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	containerName, found := c.Params.Get("containerName")
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	portNumber, found := c.Params.Get("portNumber")
	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	var workspace models.Workspace
	result := db.DB.Find(&workspace, map[string]interface{}{"ID": workspaceId, "user_id": user.ID})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	var container models.WorkspaceContainer
	result = db.DB.Find(&container, map[string]interface{}{"container_name": containerName, "workspace_id": workspace.ID})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "container not found",
		})
		return
	}

	var containerPort models.WorkspaceContainerPort
	result = db.DB.Find(&containerPort, map[string]interface{}{"container_id": container.ID, "port_number": portNumber})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "port not found",
		})
		return
	}

	db.DB.Unscoped().Delete(&containerPort)

	c.JSON(http.StatusNoContent, gin.H{
		"detail": "port has been removed",
	})
}
