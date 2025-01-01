package api

import (
	"fmt"
	"net/url"
	"strings"

	"codebox.com/api/auth"
	"codebox.com/api/cli"
	"codebox.com/api/middleware"
	"codebox.com/api/settings"
	"codebox.com/api/workspaces"
	"github.com/gin-gonic/gin"
)

func V1ApiRoutes(router *gin.Engine) {
	// middlewares
	router.Use(func(ctx *gin.Context) {
		if strings.HasPrefix(ctx.Request.Host, "codebox--") {

			splittedHostname := strings.Split(ctx.Request.Host, "--")

			if len(splittedHostname) != 4 {
				ctx.JSON(400, gin.H{
					"detail": "invalid hostname",
				})
				return
			}

			ctx.Request.URL.Path = fmt.Sprintf("/api/v1/workspace/3/container/3/forward/27017?request_path=%s", url.QueryEscape(ctx.Request.URL.Path))

			newRequestParams := []gin.Param{
				{
					Key:   "workspaceId",
					Value: "3",
				},
				{
					Key:   "containerId",
					Value: "3",
				},
				{
					Key:   "portNumber",
					Value: "27017",
				},
			}
			ctx.Params = newRequestParams

			// ctx.Request.URL.Query()

			workspaces.HandleForwardContainerPort(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
	})
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
		}

		// workspace related apis
		workspaceApis := v1.Group("/workspace")
		{
			workspaceApis.GET("", workspaces.HandleListWorkspaces)
			workspaceApis.GET("/:workspaceId", workspaces.HandleRetrieveWorkspace)
			workspaceApis.DELETE("/:workspaceId", workspaces.HandleDeleteWorkspace)
			workspaceApis.POST("", workspaces.HandleCreateWorkspace)
			workspaceApis.GET("/:workspaceId/logs", workspaces.HandleRetrieveWorkspaceLogs)
			workspaceApis.Any("/:workspaceId/container/:containerId/forward/:portNumber", workspaces.HandleForwardContainerPort)
			workspaceApis.POST("/:workspaceId/start", workspaces.HandleStartWorkspace)
			workspaceApis.POST("/:workspaceId/stop", workspaces.HandleStopWorkspace)
		}

		// instance settings related apis
		v1.GET("/instance-settings", settings.HandleRetrieveServerSettings)

		// download cli
		v1.GET("/download-cli", cli.HandleDownloadCLI)

	}
}
