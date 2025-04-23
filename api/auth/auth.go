package auth

import (
	"net/http"
	"time"

	"github.com/davidebianchi03/codebox/api/utils"
	"github.com/davidebianchi03/codebox/config"
	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

// set authentication cookie
func SetAuthCookie(ctx *gin.Context, token string) error {
	// u, err := url.Parse(config.Environment.ExternalUrl)
	// if err != nil {
	// 	return err
	// }

	// Set auth cookie, duration is set to zero because
	// token expiration has been already set in DB
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		config.Environment.AuthCookieName,
		token,
		0,
		"",
		"",
		true,
		false,
	)

	return nil
}

// POST /api/v1/login
// validate email and password and return token in response
// also set a cookie for the token
func HandleLogin(ctx *gin.Context) {
	var parsedBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := ctx.ShouldBindBodyWithJSON(&parsedBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	var user models.User
	result := db.DB.Where("email=?", parsedBody.Email).Find(&user)
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

	if !user.CheckPassword(parsedBody.Password) {
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

	result = db.DB.Create(&token)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	// Set auth cookie
	SetAuthCookie(ctx, token.Token)

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
	if err := db.DB.Model(models.User{}).Count(&usersCount).Error; err != nil {
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
	r := db.DB.Find(&users, map[string]interface{}{
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
	r = db.DB.Find(&[]models.User{}, map[string]interface{}{}).Count(&usersCount)
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
		Email:       parsedBody.Email,
		FirstName:   parsedBody.FirstName,
		LastName:    parsedBody.LastName,
		Password:    password,
		IsSuperuser: usersCount == 0,
	}

	r = db.DB.Create(&newUser)
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

	db.DB.Delete(&token)

	// clear cookies
	SetAuthCookie(ctx, "")

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
