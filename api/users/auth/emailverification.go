package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/db/models"
)

type VerifyEmailAddressRequestBody struct {
	Code string `json:"code" binding:"required"`
}

// HandleVerifyEmailAddress godoc
// @Summary Verify Email Address
// @Schemes
// @Description Verify Email Address
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body VerifyEmailAddressRequestBody true "Verification code"
// @Success 200
// @Failure 406 "Invalid verification code"
// @Failure 409 "Email already verified"
// @Router /api/v1/auth/verify-email-address [post]
func HandleVerifyEmailAddress(c *gin.Context) {
	var reqBody VerifyEmailAddressRequestBody
	err := c.ShouldBindBodyWithJSON(&reqBody)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing or invalid field")
		return
	}

	code, err := models.RetrieveVerificationCodeByCode(reqBody.Code)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	if code == nil {
		utils.ErrorResponse(
			c,
			http.StatusNotAcceptable,
			"verification code is not valid or is expired",
		)
		return
	}

	models.RevokeAllTokensForUser(code.User)

	if code.Expiration != nil {
		if time.Now().After(*code.Expiration) {
			utils.ErrorResponse(
				c,
				http.StatusNotAcceptable,
				"verification code is not valid or is expired",
			)
			return
		}
	}

	if code.User.EmailVerified {
		utils.ErrorResponse(
			c,
			http.StatusConflict,
			"email has already been verified",
		)
		return
	}

	u := code.User
	u.EmailVerified = true
	if err := models.UpdateUser(&u); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"detail": "email has been verified",
		},
	)
}
