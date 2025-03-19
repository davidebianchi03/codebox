package auth

import (
	"net/http"
	"time"

	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

// TODO: ratelimit
func HandleLogin(c *gin.Context) {
	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var parsedBody RequestBody
	err := c.ShouldBindBodyWithJSON(&parsedBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	var user models.User
	result := db.DB.Where("email=?", parsedBody.Email).Find(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"detail": "invalid credentials",
		})
		return
	}

	if !user.CheckPassword(parsedBody.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"detail": "invalid credentials",
		})
		return
	}

	token, err := models.CreateToken(user, time.Duration(time.Hour*24*20))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	result = db.DB.Create(&token)
	if result.Error != nil {
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
