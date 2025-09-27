package api

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.com/codebox4073715/codebox/api/admin"
	"gitlab.com/codebox4073715/codebox/api/auth"
	"gitlab.com/codebox4073715/codebox/api/cli"
	"gitlab.com/codebox4073715/codebox/api/middleware"
	"gitlab.com/codebox4073715/codebox/api/permissions"
	"gitlab.com/codebox4073715/codebox/api/runners"
	"gitlab.com/codebox4073715/codebox/api/settings"
	"gitlab.com/codebox4073715/codebox/api/templates"
	"gitlab.com/codebox4073715/codebox/api/workspaces"
	"gitlab.com/codebox4073715/codebox/config"
	docs "gitlab.com/codebox4073715/codebox/docs"
)

func SetupRouter() *gin.Engine {
	if config.Environment.DebugEnabled {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	V1ApiRoutes(r)
	r.LoadHTMLGlob("html/templates/*")
	return r
}

func V1ApiRoutes(router *gin.Engine) {
	// middlewares
	router.Use(middleware.PortForwardingMiddleware)
	router.Use(middleware.CORSMiddleware)

	docs.SwaggerInfo.BasePath = "/api/v1"

	// endpoints
	v1 := router.Group("/api/v1")
	{
		// auth related apis
		authApis := v1.Group("/auth")
		{
			authApis.GET(
				"/initial-user-exists",
				auth.HandleRetrieveInitialUserExists,
			)
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
			authApis.GET(
				"/subdomains/authorize",
				permissions.AuthenticationRequiredRoute(auth.HandleSubdomainLoginAuthorize),
			)
			authApis.GET(
				fmt.Sprintf("/subdomains/callback-%s", url.PathEscape(config.Environment.AuthCookieName)),
				auth.HandleSubdomainLoginCallback,
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
			workspaceApis.POST(
				"/:workspaceId/container/:containerName/port",
				permissions.AuthenticationRequiredRoute(workspaces.HandleCreateContainerPortByWorkspaceContainer),
			)
			workspaceApis.GET(
				"/:workspaceId/container/:containerName/port/:portNumber",
				permissions.AuthenticationRequiredRoute(workspaces.RetrieveContainerPortsByWorkspaceContainer),
			)
			workspaceApis.DELETE(
				"/:workspaceId/container/:containerName/port/:portNumber",
				permissions.AuthenticationRequiredRoute(workspaces.HandleDeleteContainerPortByWorkspaceContainer),
			)
			workspaceApis.Any(
				"/:workspaceId/container/:containerName/forward-http/:portNumber",
				workspaces.HandleForwardHttp,
			)
			workspaceApis.Any(
				"/:workspaceId/container/:containerName/forward-ssh",
				permissions.AuthenticationRequiredRoute(workspaces.HandleForwardSsh),
			)
			workspaceApis.Any(
				"/:workspaceId/container/:containerName/forward-tcp/:portNumber",
				permissions.AuthenticationRequiredRoute(workspaces.HandleForwardTcp),
			)
			workspaceApis.POST(
				"/:workspaceId/update-config",
				permissions.AuthenticationRequiredRoute(workspaces.HandleUpdateWorkspaceConfiguration),
			)
		}
		v1.GET("/workspace-types", permissions.AuthenticationRequiredRoute(workspaces.HandleListWorkspaceTypes))

		templatesApis := v1.Group("/templates")
		{
			templatesApis.GET("", permissions.AuthenticationRequiredRoute(templates.HandleListTemplates))
			templatesApis.GET(":templateId", permissions.AuthenticationRequiredRoute(templates.HandleRetrieveTemplate))
			templatesApis.GET(":templateId/workspaces", permissions.TemplateManagerRequiredRoute(templates.HandleListWorkspacesByTemplate))
			templatesApis.POST("", permissions.TemplateManagerRequiredRoute(templates.HandleCreateTemplate))
			templatesApis.PUT(":templateId", permissions.TemplateManagerRequiredRoute(templates.HandleUpdateTemplate))
			templatesApis.DELETE(":templateId", permissions.TemplateManagerRequiredRoute(templates.HandleDeleteWorkspace))
			templatesApis.GET(
				":templateId/versions",
				permissions.AuthenticationRequiredRoute(templates.HandleListTemplateVersionsByTemplate),
			)
			templatesApis.GET(
				":templateId/versions/:versionId",
				permissions.AuthenticationRequiredRoute(templates.HandleRetrieveTemplateVersionByTemplate),
			)
			templatesApis.GET(
				":templateId/latest-version",
				permissions.AuthenticationRequiredRoute(templates.HandleRetrieveLatestTemplateVersionByTemplate),
			)
			templatesApis.PUT(
				":templateId/versions/:versionId",
				permissions.TemplateManagerRequiredRoute(templates.HandleUpdateTemplateVersionByTemplate),
			)
			templatesApis.GET(
				":templateId/versions/:versionId/entries",
				permissions.AuthenticationRequiredRoute(templates.HandleListTemplateVersionEntries),
			)
			templatesApis.GET(
				":templateId/versions/:versionId/entries/*path",
				permissions.AuthenticationRequiredRoute(templates.HandleRetrieveTemplateVersionFile),
			)
			templatesApis.POST(
				":templateId/versions/:versionId/entries",
				permissions.TemplateManagerRequiredRoute(templates.HandleCreateTemplateVersionEntry),
			)
			templatesApis.PUT(
				":templateId/versions/:versionId/entries/*path",
				permissions.TemplateManagerRequiredRoute(templates.HandleUpdateTemplateVersionEntry),
			)
			templatesApis.DELETE(
				":templateId/versions/:versionId/entries/*path",
				permissions.TemplateManagerRequiredRoute(templates.HandleDeleteTemplateVersionEntry),
			)
		}

		templatesByName := v1.Group("/templates-by-name")
		{
			templatesByName.GET("", permissions.AuthenticationRequiredRoute(templates.HandleListTemplates))
			templatesByName.GET(":templateName", permissions.AuthenticationRequiredRoute(templates.HandleRetrieveTemplateByName))
		}

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
			permissions.AuthenticationRequiredRoute(settings.HandleRetrieveServerSettings),
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
				"stats",
				permissions.AdminRequiredRoute(admin.HandleAdminStats),
			)
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
			adminApis.DELETE(
				"users/:email",
				permissions.AdminRequiredRoute(admin.HandleAdminDeleteUser),
			)
			adminApis.POST(
				"users/:email/set-password",
				permissions.AdminRequiredRoute(admin.HandleAdminSetUserPassword),
			)
			adminApis.GET(
				"workspaces",
				permissions.AdminRequiredRoute(admin.AdminListWorkspaces),
			)
			adminApis.POST(
				"users/:email/impersonate",
				permissions.AdminRequiredRoute(admin.HandleAdminImpersonateUser),
			)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
