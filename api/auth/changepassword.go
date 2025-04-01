package auth

import (
	"net/http"

	"github.com/davidebianchi03/codebox/api/utils"
	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

// TODO: ratelimit?
func HandleChangePassword(c *gin.Context) {
	var parsedBody struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindBodyWithJSON(&parsedBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing or invalid argument",
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
		c.JSON(http.StatusExpectationFailed, gin.H{
			"detail": "invalid password",
		})
		return
	}

	// validate the new password
	if err := models.ValidatePassword(parsedBody.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
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
