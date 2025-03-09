package middleware

import (
	"net/http"

	"github.com/davidebianchi03/codebox/api/utils"
	"github.com/gin-gonic/gin"
)

func IsSuperuserMiddleware(c *gin.Context) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if !user.IsSuperuser {
		c.JSON(http.StatusForbidden, gin.H{
			"detail": "forbidden",
		})
		return
	}
	c.Next()
}
