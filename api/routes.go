package api

import (
	"codebox.com/api/auth"
	"codebox.com/api/middleware"
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
		// endpoints
		authApis := v1.Group("/auth")
		{
			authApis.POST("/login", auth.HandleLogin)
			authApis.GET("/whoami", auth.HandleWhoAmI)
		}

		workspaceApis := v1.Group("/workspace")
		{
			workspaceApis.GET("", workspaces.HandleListWorkspaces)
			workspaceApis.POST("", workspaces.HandleCreateWorkspace)
		}
	}
}
