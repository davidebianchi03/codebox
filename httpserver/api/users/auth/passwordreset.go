package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
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
	Email string `json:"email" binding:"required"`
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
