package permissions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/utils"
)

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
