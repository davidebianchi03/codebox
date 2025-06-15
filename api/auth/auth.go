package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/utils"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
)

type LoginRequestBody struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me"`
}

// Login godoc
// @Summary Login
// @Schemes
// @Description Login
// @Tags Templates
// @Accept json
// @Produce json
// @Param request body LoginRequestBody true "Credentials"
// @Success 200 {object}
// @Router /api/v1/auth/login [post]
func HandleLogin(ctx *gin.Context) {
	var requestBody *LoginRequestBody

	err := ctx.ShouldBindBodyWithJSON(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing or invalid field",
		})
		return
	}

	var user models.User
	result := dbconn.DB.Where("email=?", requestBody.Email).Find(&user)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if user.ID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"detail": "invalid credentials",
		})
		return
	}

	if !user.CheckPassword(requestBody.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"detail": "invalid credentials",
		})
		return
	}

	token, err := models.CreateToken(user, time.Duration(time.Hour*24*20))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	result = dbconn.DB.Create(&token)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	// Set auth cookie
	cookieDuration := 0
	if requestBody.RememberMe {
		cookieDuration = 3600 * 24 * 20
	}

	SetAuthCookie(ctx, token.Token, cookieDuration)

	ctx.JSON(http.StatusOK, gin.H{
		"token":      token.Token,
		"expiration": token.ExpirationDate,
	})
}

// TODO: ratelimit
// POST /api/v1/signup
// check if a user already exists,
// check if another user with the same email address exists
// validate password
func HandleSignup(ctx *gin.Context) {
	var usersCount int64
	if err := dbconn.DB.Model(models.User{}).Count(&usersCount).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if usersCount > 0 {
		ctx.JSON(http.StatusNotAcceptable, gin.H{
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

	err := ctx.ShouldBindBodyWithJSON(&parsedBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	// validate password
	if err := models.ValidatePassword(parsedBody.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	// check if user with the same email already exists
	users := []models.User{}
	r := dbconn.DB.Find(&users, map[string]interface{}{
		"email": parsedBody.Email,
	})

	if r.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if len(users) > 0 {
		ctx.JSON(http.StatusConflict, gin.H{
			"detail": "another user with the same email already exists",
		})
		return
	}

	// check the number of existing users (first user is always an admin)
	r = dbconn.DB.Find(&[]models.User{}, map[string]interface{}{}).Count(&usersCount)
	if r.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	password, err := models.HashPassword(parsedBody.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	// create new user
	newUser := models.User{
		Email:             parsedBody.Email,
		FirstName:         parsedBody.FirstName,
		LastName:          parsedBody.LastName,
		Password:          password,
		IsSuperuser:       usersCount == 0,
		IsTemplateManager: usersCount == 0,
	}

	r = dbconn.DB.Create(&newUser)
	if r.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, newUser)
}

// /api/v1/logout
// delete token and clear cookies
func HandleLogout(ctx *gin.Context) {
	token, err := utils.GetTokenFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"detail": err.Error(),
		})
	}

	dbconn.DB.Unscoped().Delete(&token)

	// clear cookies
	SetAuthCookie(ctx, "", 0)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
