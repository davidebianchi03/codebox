package permissions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
)

/*
Wrap a Gin handler to require that the user is authenticated.
If the user is not authenticated, returns 401 Unauthorized.
Otherwise, calls the original handler.
*/
func AuthenticationRequiredRoute(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := utils.GetUserFromContext(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": err.Error(),
			})
		} else {
			handler(c)
		}
	}
}
