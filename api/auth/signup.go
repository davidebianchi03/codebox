package auth

import (
	"net/http"

	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

// TODO: ratelimit
func HandleSignup(c *gin.Context) {
	var usersCount int64
	if err := db.DB.Model(models.User{}).Count(&usersCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if usersCount > 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"detail": "initial user already exists",
		})
		return
	}

	var parsedBody struct {
		Email     string `json:"email" binding:"required,email"`
		FirstName string `json:"first_name"  binding:"required"`
		LastName  string `json:"last_name"  binding:"required"`
		Password  string `json:"password"  binding:"required"`
	}

	err := c.ShouldBindBodyWithJSON(&parsedBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	// validate password
	if err := models.ValidatePassword(parsedBody.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	// check if user with the same email already exists
	users := []models.User{}
	r := db.DB.Find(&users, map[string]interface{}{
		"email": parsedBody.Email,
	})

	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if len(users) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"detail": "another user with the same email already exists",
		})
		return
	}

	// check the number of existing users (first user is always an admin)
	r = db.DB.Find(&[]models.User{}, map[string]interface{}{}).Count(&usersCount)
	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	password, err := models.HashPassword(parsedBody.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	// create new user
	newUser := models.User{
		Email:       parsedBody.Email,
		FirstName:   parsedBody.FirstName,
		LastName:    parsedBody.LastName,
		Password:    password,
		IsSuperuser: usersCount == 0,
	}

	r = db.DB.Create(&newUser)
	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}
