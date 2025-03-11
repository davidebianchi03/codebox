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
