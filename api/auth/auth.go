package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/serializers"
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
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequestBody true "Credentials"
// @Success 200 {object} serializers.TokenSerializer
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

	user, err := models.RetrieveUserByEmail(requestBody.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid credentials",
		})
		return
	}

	if !user.CheckPassword(requestBody.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid credentials",
		})
		return
	}

	token, err := models.CreateToken(*user, time.Duration(time.Hour*24*20))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	// Set auth cookie,
	// duration = 0 means that cookie expires when browser session ends
	cookieDuration := 0
	if requestBody.RememberMe {
		cookieDuration = 3600 * 24 * 20
	}

	SetAuthCookie(ctx, token.Token, cookieDuration)

	ctx.JSON(http.StatusOK, serializers.LoadTokenSerializer(&token))
}

type SignUpRequestBody struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name"  binding:"required"`
	LastName  string `json:"last_name"  binding:"required"`
	Password  string `json:"password"  binding:"required"`
}

// TODO: ratelimit
// Signup godoc
// @Summary Signup
// @Schemes
// @Description Signup
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body SignUpRequestBody true "Credentials"
// @Success 200 {object} serializers.UserSerializer
// @Router /api/v1/auth/signup [post]
func HandleSignup(ctx *gin.Context) {
	usersCount, err := models.CountAllUsers()

	if err != nil {
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

	var requestBody *SignUpRequestBody
	err = ctx.ShouldBindBodyWithJSON(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	// validate password
	if err := models.ValidatePassword(requestBody.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	// check if user with the same email already exists
	users := []models.User{}
	r := dbconn.DB.Find(&users, map[string]interface{}{
		"email": requestBody.Email,
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

	newUser, err := models.CreateUser(
		requestBody.Email,
		requestBody.FirstName,
		requestBody.LastName,
		requestBody.Password,
		usersCount == 0,
		usersCount == 0,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	ctx.JSON(
		http.StatusCreated,
		serializers.LoadCurrentUserSerializer(newUser, false),
	)
}

// Logout godoc
// @Summary Logout
// @Schemes
// @Description Logout
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 ""
// @Router /api/v1/auth/logout [post]
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
