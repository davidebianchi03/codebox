package auth

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/emails"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"

	httperrors "gitlab.com/codebox4073715/codebox/httpserver/errors"
)

// TODO: ratelimit
// GET /api/v1/auth/initial-user-exists
// retrieve if at least one user exists for the current instance of codebox
// this api is used to redirect users to signup page to create the first user
func HandleRetrieveInitialUserExists(c *gin.Context) {
	var usersCount int64
	if err := dbconn.DB.Model(models.User{}).Count(&usersCount).Error; err != nil {
		httperrors.RenderError(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exists": usersCount > 0,
	})
}

type SignUpRequestBody struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name"  binding:"required"`
	LastName  string `json:"last_name"  binding:"required"`
	Password  string `json:"password"  binding:"required"`
}

// Signup godoc
// @Summary Signup
// @Schemes
// @Description Signup
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body SignUpRequestBody true "Credentials"
// @Success 200
// @Failure 429 "Ratelimit exceeded"
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

	instanceSettings, err := models.GetSingletonModelInstance[models.InstanceSettings]()

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
	instanceSettings, err := models.GetSingletonModelInstance[models.InstanceSettings]()

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
