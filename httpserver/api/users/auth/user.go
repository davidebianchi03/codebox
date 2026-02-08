package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
)

// HandleRetriveUserDetails godoc
// @Summary Retrieve user details
// @Schemes
// @Description Retrieve details of the currently authenticated user.
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} serializers.CurrentUserSerializer
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /api/v1/auth/user-details [get]
func HandleRetriveUserDetails(c *gin.Context) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	impersonated, err := UserIsBeingImpersonated(c)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	c.JSON(
		http.StatusOK,
		serializers.LoadCurrentUserSerializer(&user, impersonated),
	)
}

type HandleUpdateUserDetailsRequestBody struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// HandleUpdateUserDetails godoc
// @Summary Update user details
// @Schemes
// @Description Update details of the currently authenticated user.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body HandleUpdateUserDetailsRequestBody true "Request Data"
// @Success 200 {object} serializers.CurrentUserSerializer
// @Failure 400 "Missing or invalid field"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /api/v1/auth/user-details [put]
func HandleUpdateUserDetails(c *gin.Context) {
	var requestBody *HandleUpdateUserDetailsRequestBody

	err := c.ShouldBindBodyWithJSON(&requestBody)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"missing or invalid field",
		)
		return
	}

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	user.FirstName = requestBody.FirstName
	user.LastName = requestBody.LastName
	dbconn.DB.Save(&user)

	c.JSON(
		http.StatusOK,
		serializers.LoadCurrentUserSerializer(&user, false),
	)
}

// HandleRetrieveUserPublicKey godoc
// @Summary Retrieve user SSH public key
// @Schemes
// @Description Retrieve user SSH public key
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} serializers.UserSshPublicKeySerializer
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /api/v1/auth/user-ssh-public-key [get]
func HandleRetrieveUserPublicKey(c *gin.Context) {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	c.JSON(
		http.StatusOK,
		serializers.LoadUserSshPublicKeySerializer(user.SshPublicKey),
	)
}

type ChangePasswordRequestBody struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// HandleChangePassword godoc
// @Summary Change user password
// @Schemes
// @Description Change password of the currently authenticated user.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequestBody true "Request Data"
// @Success 200 "Password changed successfully"
// @Failure 400 "Missing or invalid field"
// @Failure 401 "Unauthorized"
// @Failure 417 "Invalid current password"
// @Failure 500 "Internal server error"
// @Router /api/v1/auth/change-password [put]
func HandleChangePassword(c *gin.Context) {
	var parsedBody *ChangePasswordRequestBody

	if err := c.ShouldBindBodyWithJSON(&parsedBody); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"missing or invalid field",
		)
		return
	}

	user, err := utils.GetUserFromContext(c)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	// check if current password is correct
	// the current password is used to check that the user that
	// is trying to perform this operation is the owner of the account
	if !user.CheckPassword(parsedBody.CurrentPassword) {
		utils.ErrorResponse(
			c,
			http.StatusExpectationFailed,
			"invalid current password",
		)
		return
	}

	// validate the new password
	// passwordmust be at least 10 characters long and
	// include at least one uppercase letter and one special symbol (!_-,.?!)
	if err := models.ValidatePassword(parsedBody.NewPassword); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	// hash password and store it into db
	user.Password, err = models.HashPassword(parsedBody.NewPassword)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	dbconn.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"detail": "password changed",
	})
}
