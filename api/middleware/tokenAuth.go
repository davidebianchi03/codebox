package middleware

import (
	"net/http"

	"codebox.com/api/utils"
	"github.com/gin-gonic/gin"
)

var AuthNotRequiredEndpoits = [...]string{"/api/v1/auth/login", ""}

func isAuthRequired(endpoint string) bool {
	for _, noAuthEndpoint := range AuthNotRequiredEndpoits {
		if noAuthEndpoint == endpoint {
			return false
		}
	}
	return true
}

func TokenAuthMiddleware(ctx *gin.Context) {
	if !isAuthRequired(ctx.Request.URL.Path) {
		ctx.Next()
		return
	}

	user, err := utils.GetUserFromContext(ctx)
	_ = user
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail": err.Error(),
		})
		return
	}
	ctx.Next()
}
