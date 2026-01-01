package permissions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
)

/*
Wrap a Gin handler to require that the user is a template manager.
If the user is not authenticated, returns 401 Unauthorized.
If the user is authenticated but not a template manager, returns 403 Forbidden.
Otherwise, calls the original handler.
*/
func TemplateManagerRequiredRoute(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := utils.GetUserFromContext(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": err.Error(),
			})
		} else {
			if user.IsTemplateManager || user.IsSuperuser {
				handler(c)
			} else {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"detail": "forbidden",
				})
			}
		}
	}
}
