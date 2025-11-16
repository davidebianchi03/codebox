package workspaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/config"
)

// HandleListWorkspaceTypes godoc
// @Summary List workspace types
// @Schemes
// @Description List workspace types
// @Tags Workspaces
// @Accept json
// @Produce json
// @Success 200 {object} serializers.WorkspaceTypeSerializer[]
// @Router /api/v1/workspace-types [get]
func HandleListWorkspaceTypes(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		serializers.LoadMultipleWorkspaceTypeSerializer(config.ListWorkspaceTypes()),
	)
}
