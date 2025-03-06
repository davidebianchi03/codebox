package workspaces

import (
	"net/http"

	"github.com/davidebianchi03/codebox/config"
	"github.com/gin-gonic/gin"
)

func HandleListWorkspaceTypes(c *gin.Context) {
	c.JSON(http.StatusOK, config.ListWorkspaceTypes())
}
