package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/serializers"
	"gitlab.com/codebox4073715/codebox/api/utils"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
)

// TODO: ratelimit
// GET /api/v1/auth/initial-user-exists
// retrieve if at least one user exists for the current instance of codebox
// this api is used to redirect users to signup page to create the first user
func HandleRetrieveInitialUserExists(ctx *gin.Context) {
	var usersCount int64
	if err := dbconn.DB.Model(models.User{}).Count(&usersCount).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"exists": usersCount > 0,
	})
}

// GET /api/v1/auth/user-details
// retrieve details about the current user
func HandleRetriveUserDetails(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, serializers.LoadCurrentUserSerializer(&user, false))
}

// PUT or PATCH /api/v1/auth/user-details
// update user first and last name
func HandleUpdateUserDetails(ctx *gin.Context) {
	var requestBody struct {
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
	}

	err := ctx.ShouldBindBodyWithJSON(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if requestBody.FirstName != nil {
		user.FirstName = *requestBody.FirstName
	}
	if requestBody.LastName != nil {
		user.LastName = *requestBody.LastName
	}

	dbconn.DB.Save(&user)

	ctx.JSON(http.StatusOK, serializers.LoadCurrentUserSerializer(&user, false))
}

// GET /api/v1/auth/user-ssh-public-key
// retrieve user's ssh public key
func HandleRetrieveUserPublicKey(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"public_key": user.SshPublicKey,
	})
}

// TODO: ratelimit?
// POST /api/v1/auth/change-password
// change the password of current user
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

	// check if current password is correct
	// the current password is used to check that the user that
	// is trying to perform this operation is the owner of the account
	if !user.CheckPassword(parsedBody.CurrentPassword) {
		c.JSON(http.StatusExpectationFailed, gin.H{
			"detail": "invalid password",
		})
		return
	}

	// validate the new password
	// passwordmust be at least 10 characters long and
	// include at least one uppercase letter and one special symbol (!_-,.?!)
	if err := models.ValidatePassword(parsedBody.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	// hash password and store it into db
	user.Password, err = models.HashPassword(parsedBody.NewPassword)
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
