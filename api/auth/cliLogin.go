package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/utils"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
)

// TODO: replace with api keys
func HandleCliLogin(c *gin.Context) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if user.Email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"detail": "missing or invalid authorization token",
		})
		return
	}

	token, err := models.CreateToken(user, time.Duration(time.Hour*24*90))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if err = dbconn.DB.Create(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token.Token,
		"expiration": token.ExpirationDate,
	})
}
