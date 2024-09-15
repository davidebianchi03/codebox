package api

import (
	"codebox.com/api/auth"
	"codebox.com/api/middleware"
	"codebox.com/api/workspaces"
	"github.com/gin-gonic/gin"
)

func V1ApiRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		authApis := v1.Group("/auth")
		{
			authApis.POST("/login", auth.HandleLogin)
		}

		workspaceApis := v1.Group("/workspace")
		{
			workspaceApis.Use(middleware.TokenAuthMiddleware)
			workspaceApis.GET("/", workspaces.HandleListWorkspaces)
			workspaceApis.POST("/", workspaces.HandleCreateWorkspace)
		}
	}
}
