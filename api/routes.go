package api

import (
	"github.com/davidebianchi03/codebox/api/admin"
	"github.com/davidebianchi03/codebox/api/auth"
	"github.com/davidebianchi03/codebox/api/cli"
	"github.com/davidebianchi03/codebox/api/middleware"
	"github.com/davidebianchi03/codebox/api/permissions"
	"github.com/davidebianchi03/codebox/api/runners"
	"github.com/davidebianchi03/codebox/api/settings"
	"github.com/davidebianchi03/codebox/api/workspaces"
	"github.com/davidebianchi03/codebox/config"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	if config.Environment.DebugEnabled {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	V1ApiRoutes(r)
	return r
}

func V1ApiRoutes(router *gin.Engine) {
	// middlewares
	router.Use(middleware.PortForwardingMiddleware)
	router.Use(middleware.CORSMiddleware)

	// endpoints
	v1 := router.Group("/api/v1")
	{
		// auth related apis
		authApis := v1.Group("/auth")
		{
			authApis.POST(
				"/login",
				auth.HandleLogin,
			)
			authApis.POST(
				"/logout",
				auth.HandleLogout,
			)
			authApis.POST(
				"/signup",
				auth.HandleSignup,
			)
			authApis.GET(
				"/user-details",
				permissions.AuthenticationRequiredRoute(auth.HandleRetriveUserDetails),
			)
			authApis.PUT(
				"/user-details",
				permissions.AuthenticationRequiredRoute(auth.HandleUpdateUserDetails),
			)
			authApis.PATCH(
				"/user-details",
				permissions.AuthenticationRequiredRoute(auth.HandleUpdateUserDetails),
			)
			authApis.GET(
				"/user-ssh-public-key",
				permissions.AuthenticationRequiredRoute(auth.HandleRetrieveUserPublicKey),
			)
			authApis.POST(
				"/change-password",
				permissions.AuthenticationRequiredRoute(auth.HandleChangePassword),
			)
			authApis.POST(
				"/cli-login",
				permissions.AuthenticationRequiredRoute(auth.HandleCliLogin),
			)
		}

		// workspace related apis
		workspaceApis := v1.Group("/workspace")
		{
			workspaceApis.GET(
				"",
				permissions.AuthenticationRequiredRoute(workspaces.HandleListWorkspaces),
			)
			workspaceApis.GET(
				"/:workspaceId",
				permissions.AuthenticationRequiredRoute(workspaces.HandleRetrieveWorkspace),
			)
			workspaceApis.POST(
				"",
				permissions.AuthenticationRequiredRoute(workspaces.HandleCreateWorkspace),
			)
			workspaceApis.PUT(
				"/:workspaceId",
				permissions.AuthenticationRequiredRoute(workspaces.HandleUpdateWorkspace),
			)
			workspaceApis.PATCH(
				"/:workspaceId",
				permissions.AuthenticationRequiredRoute(workspaces.HandleUpdateWorkspace),
			)
			workspaceApis.DELETE(
				"/:workspaceId",
				permissions.AuthenticationRequiredRoute(workspaces.HandleDeleteWorkspace),
			)
			workspaceApis.GET(
				"/:workspaceId/logs",
				permissions.AuthenticationRequiredRoute(workspaces.HandleRetrieveWorkspaceLogs),
			)
			workspaceApis.POST(
				"/:workspaceId/start",
				permissions.AuthenticationRequiredRoute(workspaces.HandleStartWorkspace),
			)
			workspaceApis.POST(
				"/:workspaceId/stop",
				permissions.AuthenticationRequiredRoute(workspaces.HandleStopWorkspace),
			)
			workspaceApis.GET(
				"/:workspaceId/container",
				permissions.AuthenticationRequiredRoute(workspaces.ListWorkspaceContainersByWorkspace),
			)
			workspaceApis.GET(
				"/:workspaceId/container/:containerName",
				permissions.AuthenticationRequiredRoute(workspaces.RetrieveWorkspaceContainersByWorkspace),
			)
			workspaceApis.GET(
				"/:workspaceId/container/:containerName/port",
				permissions.AuthenticationRequiredRoute(workspaces.ListContainerPortsByWorkspaceContainer),
			)
			workspaceApis.GET(
				"/:workspaceId/container/:containerName/port/:portNumber",
				permissions.AuthenticationRequiredRoute(workspaces.RetrieveContainerPortsByWorkspaceContainer),
			)
			workspaceApis.Any(
				"/:workspaceId/container/:containerName/forward-http/:portNumber",
				permissions.AuthenticationRequiredRoute(workspaces.HandleForwardHttp),
			)
			workspaceApis.Any(
				"/:workspaceId/container/:containerName/forward-ssh",
				permissions.AuthenticationRequiredRoute(workspaces.HandleForwardSsh),
			)
			workspaceApis.POST(
				"/:workspaceId/update-config",
				permissions.AuthenticationRequiredRoute(workspaces.HandleUpdateWorkspaceConfiguration),
			)
		}
		v1.GET("/workspace-types", permissions.AuthenticationRequiredRoute(workspaces.HandleListWorkspaceTypes))

		// runners related apis
		runnersApis := v1.Group("/runners")
		{
			runnersApis.GET("", permissions.AuthenticationRequiredRoute(runners.HandleListRunners))
			runnersApis.Any(":runnerId/connect", runners.HandleRunnerConnect)
		}
		v1.GET(
			"/runner-types",
			permissions.AuthenticationRequiredRoute(runners.HandleListRunnerTypes),
		)

		// instance settings related apis
		v1.GET(
			"/instance-settings",
			settings.HandleRetrieveServerSettings,
		)

		// download cli
		v1.GET(
			"/download-cli",
			cli.HandleDownloadCLI,
		)

		// admin routes
		adminApis := v1.Group("/admin")
		{
			adminApis.GET(
				"runners",
				permissions.AdminRequiredRoute(admin.HandleAdminListRunners),
			)
			adminApis.GET(
				"runners/:runnerId",
				permissions.AdminRequiredRoute(admin.HandleAdminRetrieveRunners),
			)
			adminApis.PUT(
				"runners/:runnerId",
				permissions.AdminRequiredRoute(admin.HandleAdminUpdateRunner),
			)
			adminApis.POST(
				"runners",
				permissions.AdminRequiredRoute(admin.HandleAdminCreateRunner),
			)
			adminApis.GET(
				"users",
				permissions.AdminRequiredRoute(admin.HandleAdminListUsers),
			)
			adminApis.POST(
				"users",
				permissions.AdminRequiredRoute(admin.HandleAdminCreateUser),
			)
			adminApis.GET(
				"users/:email",
				permissions.AdminRequiredRoute(admin.HandleAdminRetrieveUser),
			)
			adminApis.PUT(
				"users/:email",
				permissions.AdminRequiredRoute(admin.HandleAdminUpdateUser),
			)
			adminApis.PATCH(
				"users/:email",
				permissions.AdminRequiredRoute(admin.HandleAdminUpdateUser),
			)
			adminApis.POST(
				"users/:email/set-password",
				permissions.AdminRequiredRoute(admin.HandleAdminSetUserPassword),
			)
		}
	}
}
