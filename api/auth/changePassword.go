package auth

import (
	"net/http"

	"github.com/davidebianchi03/codebox/api/utils"
	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

func HandleChangePassword(c *gin.Context) {
	type RequestBody struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	var parsedBody RequestBody
	err := c.ShouldBindBodyWithJSON(&parsedBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if !user.CheckPassword(parsedBody.CurrentPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"detail": "invalid credentials",
		})
		return
	}

	user.Password, err = models.HashPassword(parsedBody.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	db.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"detail": "password changed",
	})
}
