package permissions

import (
	"net/http"

	"github.com/davidebianchi03/codebox/api/utils"
	"github.com/gin-gonic/gin"
)

func AdminRequiredRoute(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := utils.GetUserFromContext(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": err.Error(),
			})
		} else {
			if user.IsSuperuser {
				handler(c)
			} else {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"detail": "forbidden",
				})
			}
		}
	}
}
