package api

import (
	"codebox.com/api/auth"
	"github.com/gin-gonic/gin"
)

func V1ApiRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		authApis := v1.Group("/auth")
		{
			authApis.POST("/login", auth.HandleLogin)
		}
	}
}
