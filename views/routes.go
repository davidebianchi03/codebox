package views

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/permissions"
)

func ViewsRoutes(router *gin.Engine) {
	viewsRoutes := router.Group("/views")
	{
		viewsRoutes.GET(
			"/workspace/:workspaceId/container/:containerName/terminal",
			permissions.AuthenticationRequiredRoute(HandleTerminalView),
		)
		viewsRoutes.Any(
			"/port-forward/workspace/:workspaceId/container/:containerName/port/:portNumber",
			HandlePortForwardView,
		)
	}
}
