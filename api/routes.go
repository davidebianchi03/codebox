package api

import (
	"github.com/davidebianchi03/codebox/api/auth"
	"github.com/davidebianchi03/codebox/api/middleware"
	"github.com/davidebianchi03/codebox/api/workspaces"
	"github.com/gin-gonic/gin"
)

func V1ApiRoutes(router *gin.Engine) {
	// middlewares
	router.Use(middleware.PortForwardingMiddleware)
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.TokenAuthMiddleware)

	// endpoints
	v1 := router.Group("/api/v1")
	{
		// auth related apis
		authApis := v1.Group("/auth")
		{
			authApis.POST("/login", auth.HandleLogin)
			authApis.POST("/logout", auth.HandleLogout)
			authApis.GET("/user-details", auth.HandleRetriveUserDetails)
			authApis.PATCH("/user-details", auth.HandleUpdateUserDetails)
			authApis.POST("/change-password", auth.HandleChangePassword)
			authApis.POST("/signup", auth.HandleSignup)
		}

		// // workspace related apis
		workspaceApis := v1.Group("/workspace")
		{
			// workspaceApis.GET("", workspaces.HandleListWorkspaces)
			// workspaceApis.GET("/:workspaceId", workspaces.HandleRetrieveWorkspace)
			// workspaceApis.DELETE("/:workspaceId", workspaces.HandleDeleteWorkspace)
			workspaceApis.POST("", workspaces.HandleCreateWorkspace)
			// workspaceApis.GET("/:workspaceId/logs", workspaces.HandleRetrieveWorkspaceLogs)
			workspaceApis.Any("/:workspaceId/container/:containerName/forward-http/:portNumber", workspaces.HandleForwardHttp)
			// workspaceApis.POST("/:workspaceId/start", workspaces.HandleStartWorkspace)
			// workspaceApis.POST("/:workspaceId/stop", workspaces.HandleStopWorkspace)
		}

		// // instance settings related apis
		// v1.GET("/instance-settings", settings.HandleRetrieveServerSettings)

		// // download cli
		// v1.GET("/download-cli", cli.HandleDownloadCLI)

	}
}
