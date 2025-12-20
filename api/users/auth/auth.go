package auth

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
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
// @Success 200
// @Router /api/v1/auth/signup [post]
func HandleSignup(c *gin.Context) {
	// if user is already logged in return an error
	_, err := utils.GetUserFromContext(c)
	if err == nil {
		utils.ErrorResponse(
			c,
			http.StatusUnauthorized,
			"logout before signing up",
		)
		return
	}

	instanceSettings, err := models.GetInstanceSettings()

	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	usersCount, err := models.CountAllUsers()

	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	if usersCount > 0 && !instanceSettings.IsSignUpOpen {
		utils.ErrorResponse(
			c,
			http.StatusNotAcceptable,
			"cannot signup",
		)
		return
	}

	var requestBody *SignUpRequestBody
	err = c.ShouldBindBodyWithJSON(&requestBody)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	canSignup := false

	if usersCount > 0 {
		// check if email matches an allowed pattern if signup is restricted
		if instanceSettings.IsSignUpRestricted {
			if len(strings.TrimSpace(instanceSettings.AllowedEmailRegex)) > 0 {
				allowedEmailsRegex := strings.Split(instanceSettings.AllowedEmailRegex, "\n")
				for _, re := range allowedEmailsRegex {
					m, err := regexp.MatchString(strings.TrimSpace(re), requestBody.Email)
					if err != nil {
						utils.ErrorResponse(
							c,
							http.StatusInternalServerError,
							"internal server error",
						)
						return
					}

					if m {
						canSignup = true
					}
				}
			} else {
				canSignup = true
			}
		} else {
			canSignup = true
		}

		// check if email matches a blackisted pattern
		if len(strings.TrimSpace(instanceSettings.BlockedEmailRegex)) > 0 {
			blackistedEmailsRegex := strings.Split(instanceSettings.BlockedEmailRegex, "\n")
			for _, re := range blackistedEmailsRegex {
				m, err := regexp.MatchString(strings.TrimSpace(re), requestBody.Email)
				if err != nil {
					utils.ErrorResponse(
						c,
						http.StatusInternalServerError,
						"internal server error",
					)
					return
				}

				if m {
					canSignup = false
				}
			}
		}
	} else {
		canSignup = true
	}

	if !canSignup {
		utils.ErrorResponse(
			c,
			http.StatusNotAcceptable,
			"cannot signup",
		)
		return
	}

	// validate password
	if err := models.ValidatePassword(requestBody.Password); err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	// check if user with the same email already exists
	existingUser, err := models.RetrieveUserByEmail(requestBody.Email)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	if existingUser != nil {
		// send email notifying about existing account
		// but do not reveal that the account exists in the response
		emails.SendUserAlreadyExistsEmail(requestBody.Email)
		c.JSON(
			http.StatusCreated,
			gin.H{"detail": "account created successfully"},
		)
		return
	}

	_, err = models.CreateUser(
		requestBody.Email,
		requestBody.FirstName,
		requestBody.LastName,
		requestBody.Password,
		usersCount == 0,
		usersCount == 0,
		usersCount == 0,
	)

	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		gin.H{"detail": "account created successfully"},
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
		utils.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	if token.ImpersonatedUser != nil {
		// stop impersonation log
		log, err := models.RetrieveLatestImpersonationLogByToken(token)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
			return
		}

		now := time.Now()
		log.ImpersonationFinishedAt = &now
		if err := models.UpdateImpersonationLog(log); err != nil {
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	models.DeleteToken(&token)

	// clear cookies
	SetAuthCookie(ctx, "", 0)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// HandleIsSignUpOpen godoc
// @Summary Check if signup is open
// @Schemes
// @Description Check if signup is open
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} serializers.IsSignUpOpenSerializer
// @Router /api/v1/auth/is-signup-open [get]
func HandleIsSignUpOpen(c *gin.Context) {
	instanceSettings, err := models.GetInstanceSettings()

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
		serializers.LoadIsSignUpOpenSerializer(instanceSettings.IsSignUpOpen),
	)
}

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
