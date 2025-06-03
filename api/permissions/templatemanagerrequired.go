package permissions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/utils"
)

func TemplateManagerRequiredRoute(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := utils.GetUserFromContext(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": err.Error(),
			})
		} else {
			if user.IsTemplateManager {
				handler(c)
			} else {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"detail": "forbidden",
				})
			}
		}
	}
}
