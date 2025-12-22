package auth

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/emails"
)

type LoginRequestBody struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RememberMe bool   `json:"remember_me"`
}

// TODO: ratelimit
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
func HandleLogin(c *gin.Context) {
	_, err := utils.GetUserFromContext(c)
	if err == nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"already logged in",
		)
		return
	}

	var requestBody *LoginRequestBody

	err = c.ShouldBindBodyWithJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing or invalid field",
		})
		return
	}

	user, err := models.RetrieveUserByEmail(requestBody.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid credentials",
		})
		return
	}

	if !user.CheckPassword(requestBody.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid credentials",
		})
		return
	}

	// check if user email has been verified
	if !user.EmailVerified {
		// create verification code
		err := models.RevokeAllTokensForUser(*user)
		if err != nil {
			// TODO: log error
			c.JSON(http.StatusInternalServerError, gin.H{
				"detail": "internal server error",
			})
			return
		}

		expirationDate := time.Now().Add(time.Hour * 24)
		code, err := models.CreateEmailVerificationCode(&expirationDate, *user)
		if err != nil {
			// TODO: log error
			c.JSON(http.StatusInternalServerError, gin.H{
				"detail": "internal server error",
			})
			return
		}

		// send email verification token
		verificationUrl := config.Environment.ExternalUrl + "/verify-email?code=" + url.QueryEscape(code.Code)
		err = emails.SendEmailVerificationEmail(*user, verificationUrl)
		if err != nil {
			// TODO: log error
			c.JSON(http.StatusInternalServerError, gin.H{
				"detail": "internal server error",
			})
			return
		}

		utils.ErrorResponse(
			c,
			http.StatusPreconditionFailed,
			"the email address has not yet been verified",
		)
		return
	}

	token, err := models.CreateToken(*user, time.Duration(time.Hour*24*20))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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

	SetAuthCookie(c, token.Token, cookieDuration)

	c.JSON(http.StatusOK, serializers.LoadTokenSerializer(&token))
}
