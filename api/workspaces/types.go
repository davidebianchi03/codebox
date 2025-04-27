package workspaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
)

func HandleListWorkspaceTypes(c *gin.Context) {
	c.JSON(http.StatusOK, config.ListWorkspaceTypes())
}
