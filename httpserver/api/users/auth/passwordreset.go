package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/emails"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
)

// HandleCanResetPassword godoc
// @Summary Can Reset Password
// @Schemes
// @Description Check if password reset is available
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} serializers.CanResetPasswordSerializer
// @Failure 500 "Internal server error"
// @Router /api/v1/auth/can-reset-password [get]
func HandleCanResetPassword(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		serializers.LoadCanResetPasswordSerializer(config.IsEmailConfigured()),
	)
}

type RequestPasswordResetTokenBody struct {
	Email string `json:"email" binding:"required,email"`
}

// HandleCanResetPassword godoc
// @Summary Can Reset Password
// @Schemes
// @Description Check if password reset is available
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RequestPasswordResetTokenBody true "Request Data"
// @Success 200 {object} serializers.RequestPasswordResetSerializer
// @Failure 400 "Missing or invalid field"
// @Failure 406 "Password reset is not available"
// @Failure 429 "Rate limit exceeded"
// @Failure 500 "Internal server error"
// @Router /api/v1/auth/request-password-reset [post]
func HandleRequestPasswordReset(c *gin.Context) {
	_, err := utils.GetUserFromContext(c)
	if err == nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"already logged in",
		)
		return
	}

	if !config.IsEmailConfigured() {
		utils.ErrorResponse(
			c,
			http.StatusNotAcceptable,
			"password reset is not available",
		)
		return
	}

	var requestBody *RequestPasswordResetTokenBody
	err = c.ShouldBindBodyWithJSON(&requestBody)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"missing or invalid field",
		)
		return
	}

	user, err := models.RetrieveUserByEmail(requestBody.Email)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	if user != nil {
		// send password reset email
		// generate password reset token
		prt, err := models.CreatePasswordResetToken(*user)
		if err != nil {
			utils.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"internal server error",
			)
			return
		}

		if err := emails.SendPasswordResetEmail(
			user.Email,
			prt.Token,
		); err != nil {
			utils.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"internal server error",
			)
			return
		}
	} else {
		// send user not registered email (to avoid user enumeration)
		if err := emails.SendUserNotRegisteredEmail(
			requestBody.Email,
		); err != nil {
			utils.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"internal server error",
			)
			return
		}
	}

	c.JSON(
		http.StatusOK,
		serializers.LoadRequestPasswordResetSerializer(
			true, requestBody.Email,
		),
	)
}

type HandlePasswordResetFromTokenBody struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// HandlePasswordResetFromToken godoc
// @Summary Reset Password From Token
// @Schemes
// @Description Reset password using a token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body HandlePasswordResetFromTokenBody true "Request Data"
// @Success 200 "Password reset successfully"
// @Failure 400 "Missing or invalid field"
// @Failure 404 "Invalid or expired token"
// @Failure 406 "Password reset is not available"
// @Failure 429 "Rate limit exceeded"
// @Failure 500 "Internal server error"
// @Router /api/v1/auth/password-reset-from-token [post]
func HandlePasswordResetFromToken(c *gin.Context) {
	_, err := utils.GetUserFromContext(c)
	if err == nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"already logged in",
		)
		return
	}

	if !config.IsEmailConfigured() {
		utils.ErrorResponse(
			c,
			http.StatusNotAcceptable,
			"password reset is not available",
		)
		return
	}

	var requestBody *HandlePasswordResetFromTokenBody
	err = c.ShouldBindBodyWithJSON(&requestBody)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"missing or invalid field",
		)
		return
	}

	// delete expired tokens
	if err := models.DeleteExpiredPasswordResetTokens(); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	// get password reset token from db
	prt, err := models.GetPasswordResetToken(requestBody.Token)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	if prt == nil {
		utils.ErrorResponse(
			c,
			http.StatusNotFound,
			"invalid or expired token",
		)
		return
	}

	// validate the new password
	// passwordmust be at least 10 characters long and
	// include at least one uppercase letter and one special symbol (!_-,.?!)
	if err := models.ValidatePassword(requestBody.NewPassword); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	// hash password and store it into db
	user := prt.User
	user.Password, err = models.HashPassword(requestBody.NewPassword)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	emails.SendPasswordResetDoneEmail(user.Email)

	dbconn.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"detail": "password has been reset successfully",
	})
}
