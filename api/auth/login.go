package auth

import (
	"net/http"

	"codebox.com/db"
	"github.com/gin-gonic/gin"
)

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

	var users []db.User
	result := db.DB.Model(db.User{Email: parsedBody.Email}).Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if len(users) != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"detail": "invalid credentials",
		})
		return
	}

	user := users[0]
	if !user.CheckPassword(parsedBody.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"detail": "invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"detail": "success",
	})
}
