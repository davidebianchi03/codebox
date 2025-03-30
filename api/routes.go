package api

import (
	"github.com/davidebianchi03/codebox/api/admin"
	"github.com/davidebianchi03/codebox/api/auth"
	"github.com/davidebianchi03/codebox/api/cli"
	"github.com/davidebianchi03/codebox/api/middleware"
	"github.com/davidebianchi03/codebox/api/runners"
	"github.com/davidebianchi03/codebox/api/settings"
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
			authApis.PUT("/user-details", auth.HandleUpdateUserDetails)
			authApis.PATCH("/user-details", auth.HandleUpdateUserDetails)
			authApis.GET("/user-ssh-public-key", auth.HandleRetrieveUserPublicKey)
			authApis.POST("/change-password", auth.HandleChangePassword)
			authApis.POST("/signup", auth.HandleSignup)
			authApis.POST("/cli-login", auth.HandleCliLogin)
		}

		// workspace related apis
		workspaceApis := v1.Group("/workspace")
		{
			workspaceApis.GET("", workspaces.HandleListWorkspaces)
			workspaceApis.GET("/:workspaceId", workspaces.HandleRetrieveWorkspace)
			workspaceApis.POST("", workspaces.HandleCreateWorkspace)
			workspaceApis.DELETE("/:workspaceId", workspaces.HandleDeleteWorkspace)
			workspaceApis.GET("/:workspaceId/logs", workspaces.HandleRetrieveWorkspaceLogs)
			workspaceApis.POST("/:workspaceId/start", workspaces.HandleStartWorkspace)
			workspaceApis.POST("/:workspaceId/stop", workspaces.HandleStopWorkspace)
			workspaceApis.GET("/:workspaceId/container", workspaces.ListWorkspaceContainersByWorkspace)
			workspaceApis.GET("/:workspaceId/container/:containerName", workspaces.RetrieveWorkspaceContainersByWorkspace)
			workspaceApis.GET("/:workspaceId/container/:containerName/port", workspaces.ListContainerPortsByWorkspaceContainer)
			workspaceApis.GET("/:workspaceId/container/:containerName/port/:portNumber", workspaces.RetrieveContainerPortsByWorkspaceContainer)
			workspaceApis.Any("/:workspaceId/container/:containerName/forward-http/:portNumber", workspaces.HandleForwardHttp)
			workspaceApis.Any("/:workspaceId/container/:containerName/forward-ssh", workspaces.HandleForwardSsh)
		}
		v1.GET("/workspace-types", workspaces.HandleListWorkspaceTypes)

		// runners related apis
		runnersApis := v1.Group("/runners")
		{
			runnersApis.GET("", runners.HandleListRunners)
			runnersApis.Any(":runnerId/connect", runners.HandleRunnerConnect)
		}
		v1.GET("/runner-types", runners.HandleListRunnerTypes)

		// instance settings related apis
		v1.GET("/instance-settings", settings.HandleRetrieveServerSettings)

		// download cli
		v1.GET("/download-cli", cli.HandleDownloadCLI)

		// admin routes
		adminApis := v1.Group("/admin")
		{
			adminApis.GET("runners", admin.HandleAdminListRunners)
			adminApis.GET("runners/:runnerId", admin.HandleAdminRetrieveRunners)
			adminApis.PUT("runners/:runnerId", admin.HandleAdminUpdateRunner)
			adminApis.POST("runners", admin.HandleAdminCreateRunner)
			adminApis.GET("users", admin.HandleAdminListUsers)
			adminApis.POST("users", admin.HandleAdminCreateUser)
			adminApis.GET("users/:email", admin.HandleAdminRetrieveUser)
			adminApis.PUT("users/:email", admin.HandleAdminUpdateUser)
			adminApis.PATCH("users/:email", admin.HandleAdminUpdateUser)
			adminApis.PATCH("users/:email/set-password", admin.HandleAdminSetUserPassword)
			runnersApis.Use(middleware.IsSuperuserMiddleware)
		}
	}
}
