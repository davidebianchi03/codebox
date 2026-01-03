package auth

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/emails"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
)

// HandleRetrieveInitialUserExists godoc
// @Summary Check if at lease one user exists
// @Schemes
// @Description retrieve if at least one user exists for the current instance of codebox
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} serializers.InitialUserExistsSerializer
// @Failure 429 "Ratelimit exceeded"
// @Router /api/v1/auth/initial-user-exists [get]
func HandleRetrieveInitialUserExists(c *gin.Context) {
	count, err := models.CountAllUsers()
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
		serializers.LoadInitialUserExistsSerializer(count > 0),
	)
}

/*
Check if email is matching at least one of the given regex
*/
func IsEmailMatchingARegex(email string, regEx []string) bool {
	for _, re := range regEx {
		m, err := regexp.MatchString(strings.TrimSpace(re), email)
		if err != nil {
			return false
		}

		if m {
			return true
		}
	}

	return false
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

	s, err := models.GetSingletonModelInstance[models.AuthenticationSettings]()

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

	if usersCount > 0 && !s.IsSignUpOpen {
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
		if s.IsSignUpRestricted {
			if len(strings.TrimSpace(s.AllowedEmailRegex)) > 0 {
				canSignup = IsEmailMatchingARegex(
					requestBody.Email,
					strings.Split(s.AllowedEmailRegex, "\n"),
				)
			} else {
				canSignup = true
			}
		} else {
			canSignup = true
		}

		// check if email matches a blackisted pattern
		if len(strings.TrimSpace(s.BlockedEmailRegex)) > 0 {
			match := IsEmailMatchingARegex(
				requestBody.Email,
				strings.Split(s.BlockedEmailRegex, "\n"),
			)
			if match {
				canSignup = false
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

	// check if user is approved automatically
	autoApproved := false
	if len(strings.TrimSpace(s.ApprovedByDefaultEmailRegex)) > 0 {
		autoApproved = IsEmailMatchingARegex(
			requestBody.Email,
			strings.Split(s.ApprovedByDefaultEmailRegex, "\n"),
		)
	}

	if usersCount == 0 {
		autoApproved = true
	}

	_, err = models.CreateUser(
		requestBody.Email,
		requestBody.FirstName,
		requestBody.LastName,
		requestBody.Password,
		usersCount == 0,
		usersCount == 0,
		usersCount == 0,
		autoApproved,
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
	s, err := models.GetSingletonModelInstance[models.AuthenticationSettings]()

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
		serializers.LoadIsSignUpOpenSerializer(s.IsSignUpOpen),
	)
}
