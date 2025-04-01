package permissions

import (
	"net/http"

	"github.com/davidebianchi03/codebox/api/utils"
	"github.com/gin-gonic/gin"
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
