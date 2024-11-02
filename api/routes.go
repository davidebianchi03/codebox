package api

import (
	"codebox.com/api/auth"
	"codebox.com/api/middleware"
	"codebox.com/api/settings"
	"codebox.com/api/workspaces"
	"github.com/gin-gonic/gin"
)

func V1ApiRoutes(router *gin.Engine) {

	// middlewares
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.TokenAuthMiddleware)

	// endpoints
	v1 := router.Group("/api/v1")
	{
		// auth related apis
		authApis := v1.Group("/auth")
		{
			authApis.POST("/login", auth.HandleLogin)
			authApis.GET("/user-details", auth.HandleRetriveUserDetails)
		}

		// workspace related apis
		workspaceApis := v1.Group("/workspace")
		{
			workspaceApis.GET("", workspaces.HandleListWorkspaces)
			workspaceApis.GET("/:workspaceId", workspaces.HandleRetrieveWorkspace)
			workspaceApis.POST("", workspaces.HandleCreateWorkspace)
			workspaceApis.GET("/:workspaceId/logs", workspaces.HandleRetrieveWorkspaceLogs)
			workspaceApis.Any("/:workspaceId/container/:containerId/forward", workspaces.HandleForwardContainerPort)
			workspaceApis.POST("/:workspaceId/start", workspaces.HandleStartWorkspace)
			workspaceApis.POST("/:workspaceId/stop", workspaces.HandleStopWorkspace)
		}

		// instance settings related apis
		v1.GET("/instance-settings", settings.HandleRetrieveServerSettings)
	}
}
