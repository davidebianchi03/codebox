package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
)

func HandleAdminListUsers(c *gin.Context) {
	var users *[]models.User

	if dbconn.DB.Find(&users).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

func HandleAdminRetrieveUser(c *gin.Context) {
	var user *models.User
	email, _ := c.Params.Get("email")

	if dbconn.DB.Find(&user, map[string]interface{}{
		"email": email,
	}).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func HandleAdminCreateUser(c *gin.Context) {
	var reqBody struct {
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required"`
		FirstName   string `json:"first_name" binding:"required"`
		LastName    string `json:"last_name" binding:"required"`
		IsSuperuser bool   `json:"is_superuser"`
	}

	if c.ShouldBindBodyWithJSON(&reqBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing or invalid field",
		})
		return
	}

	// check if exists another user with the same email address
	users := []models.User{}
	r := dbconn.DB.Find(&users, map[string]interface{}{
		"email": reqBody.Email,
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

	// validate password
	if err := models.ValidatePassword(reqBody.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	password, err := models.HashPassword(reqBody.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	// create new user
	newUser := models.User{
		Email:       reqBody.Email,
		FirstName:   reqBody.FirstName,
		LastName:    reqBody.LastName,
		Password:    password,
		IsSuperuser: reqBody.IsSuperuser,
	}

	r = dbconn.DB.Create(&newUser)
	if r.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func HandleAdminUpdateUser(c *gin.Context) {
	var user *models.User
	email, _ := c.Params.Get("email")

	if dbconn.DB.Find(&user, map[string]interface{}{
		"email": email,
	}).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "user not found",
		})
		return
	}

	var requestBody struct {
		FirstName   *string `json:"first_name"`
		LastName    *string `json:"last_name"`
		IsSuperuser *bool   `json:"is_superuser"`
	}

	if c.ShouldBindBodyWithJSON(&requestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid or missing argument",
		})
		return
	}

	if requestBody.FirstName != nil {
		user.FirstName = *requestBody.FirstName
	}

	if requestBody.LastName != nil {
		user.LastName = *requestBody.LastName
	}

	if requestBody.IsSuperuser != nil {
		user.IsSuperuser = *requestBody.IsSuperuser
	}

	dbconn.DB.Save(&user)

	c.JSON(http.StatusOK, user)
}

func HandleAdminSetUserPassword(c *gin.Context) {
	var user *models.User
	var err error
	email, _ := c.Params.Get("email")

	if dbconn.DB.Find(&user, map[string]interface{}{
		"email": email,
	}).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"detail": "user not found",
		})
		return
	}

	var requestBody struct {
		Password string `json:"password" binding:"required"`
	}

	if c.ShouldBindBodyWithJSON(&requestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid or missing argument",
		})
		return
	}

	// validate password
	if err := models.ValidatePassword(requestBody.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	user.Password, err = models.HashPassword(requestBody.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	dbconn.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"detail": "password changed",
	})
}
