package api

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.com/codebox4073715/codebox/config"
	docs "gitlab.com/codebox4073715/codebox/docs"
	runnerapis "gitlab.com/codebox4073715/codebox/httpserver/api/runner"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/admin"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/admin/email"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/admin/settings"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/auth"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/cli"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/common"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/runners"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/templates"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/workspaces"
	"gitlab.com/codebox4073715/codebox/httpserver/permissions"
)

func V1ApiRoutes(router *gin.Engine) {

	docs.SwaggerInfo.BasePath = "/api/v1"

	// endpoints
	v1 := router.Group("/api/v1")
	{
		// common
		v1.GET("/version", common.HandleRetrieveServerVersion)

		// auth related apis
		authApis := v1.Group("/auth")
		{
			authApis.GET(
				"/initial-user-exists",
				auth.HandleRetrieveInitialUserExists,
			)
			authApis.POST(
				"/login",
				permissions.IPRateLimitedRoute(
					auth.HandleLogin,
					8,
					60,
				),
			)
			authApis.POST(
				"/logout",
				permissions.AuthenticationRequiredRoute(auth.HandleLogout),
			)
			authApis.POST(
				"/signup",
				permissions.IPRateLimitedRoute(
					auth.HandleSignup,
					5,
					60,
				),
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
			authApis.GET(
				"/is-signup-open",
				auth.HandleIsSignUpOpen,
			)
			authApis.POST(
				"/verify-email-address",
				auth.HandleVerifyEmailAddress,
			)
			authApis.GET(
				"/can-reset-password",
				auth.HandleCanResetPassword,
			)
			authApis.POST(
				"/request-password-reset",
				permissions.IPRateLimitedRoute(
					auth.HandleRequestPasswordReset,
					5,
					60,
				),
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
				"/:workspaceId/container/:containerName/forward-ssh",
				permissions.AuthenticationRequiredRoute(workspaces.HandleForwardSsh),
			)
			workspaceApis.Any(
				"/:workspaceId/container/:containerName/terminal",
				permissions.AuthenticationRequiredRoute(workspaces.HandleTerminal),
			)
			workspaceApis.POST(
				"/:workspaceId/update-config",
				permissions.AuthenticationRequiredRoute(workspaces.HandleUpdateWorkspaceConfiguration),
			)
			workspaceApis.POST(
				"/:workspaceId/set-runner",
				permissions.AuthenticationRequiredRoute(workspaces.HandleSetRunnerForWorkspace),
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
		}
		v1.GET(
			"/runner-types",
			permissions.AuthenticationRequiredRoute(runners.HandleListRunnerTypes),
		)

		// cli
		v1.GET("/cli-version", cli.HandleRetrieveCLIVersion)
		v1.GET("/cli", cli.HandleListCLI)
		v1.GET("/cli/:id", cli.HandleRetrieveCLI)
		v1.GET(
			"/cli/:id/download",
			cli.HandleDownloadCLI,
		)

		v1.POST("stop-impersonation", permissions.AuthenticationRequiredRoute(auth.HandleStopImpersonation))

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
			adminApis.POST(
				"runners",
				permissions.AdminRequiredRoute(admin.HandleAdminCreateRunner),
			)
			adminApis.PUT(
				"runners/:runnerId",
				permissions.AdminRequiredRoute(admin.HandleAdminUpdateRunner),
			)
			adminApis.DELETE(
				"runners/:runnerId",
				permissions.AdminRequiredRoute(admin.HandleAdminDeleteRunner),
			)
			adminApis.GET(
				"recommended-runner-version",
				permissions.AdminRequiredRoute(admin.HandleRetrieveRecommendedRunnerVersion),
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
			adminApis.POST(
				"users/:email/impersonate",
				permissions.AdminRequiredRoute(admin.HandleAdminImpersonateUser),
			)
			adminApis.GET(
				"users/:email/impersonation-logs",
				permissions.AdminRequiredRoute(admin.HandleAdminListImpersonationLogsByUser),
			)
			// instance settings related apis
			adminApis.GET(
				"authentication-settings",
				permissions.AdminRequiredRoute(settings.HandleRetrieveAuthenticationSettings),
			)
			adminApis.PUT(
				"authentication-settings",
				permissions.AdminRequiredRoute(settings.HandleUpdateAuthenticationSettings),
			)
			adminApis.GET(
				"email-service-configured",
				permissions.AdminRequiredRoute(common.HandleAdminEmailServiceConfigured),
			)
			adminApis.POST(
				"send-test-email",
				permissions.AdminRequiredRoute(email.HandleSendTestEmail),
			)
		}
	}

	runnerAPIGroup := router.Group("/runner-api/v1/")
	{
		runnerAPIGroup.POST(
			"runners/:runnerId/request-port",
			permissions.RunnerTokenAuthenticationRequired(runnerapis.HandleRunnerRequestPort),
		)
		runnerAPIGroup.Any(
			"runners/:runnerId/connect",
			permissions.RunnerTokenAuthenticationRequired(runnerapis.HandleRunnerConnect),
		)
		runnerAPIGroup.Any(
			"runners/:runnerId/workspaces/:workspaceId/container/:containerName/git-ssh",
			permissions.RunnerTokenAuthenticationRequired(runnerapis.HandleRunnerGitSSH),
		)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
