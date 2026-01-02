package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
)

// HandleRetrieveAuthenticationSettings godoc
// @Summary Retrieve authentication settings
// @Schemes
// @Description Retrieve authentication settings, this api is available only to administrators
// @Tags Templates
// @Accept json
// @Produce json
// @Success 200 {object} serializers.AuthenticationSettingsSerializer
// @Router /api/v1/admin/authentication-settings [get]
func HandleRetrieveAuthenticationSettings(c *gin.Context) {
	is, err := models.GetSingletonModelInstance[models.AuthenticationSettings]()
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	c.JSON(http.StatusOK, serializers.LoadAuthenticationSettingsSerializer(is))
}

type HandleUpdateServerSettingsRequestBody struct {
	IsSignUpOpen                *bool   `json:"is_signup_open" binding:"required"`
	IsSignUpRestricted          *bool   `json:"is_signup_restricted" binding:"required"`
	AllowedEmailRegex           *string `json:"allowed_emails_regex" binding:"required"`
	BlockedEmailRegex           *string `json:"blocked_emails_regex" binding:"required"`
	UsersMustBeApproved         *bool   `json:"users_must_be_approved" binding:"required"`
	ApprovedByDefaultEmailRegex *string `json:"approved_by_default_emails_regex" binding:"required"`
}

// HandleUpdateAuthenticationSettings godoc
// @Summary Update authentication settings
// @Schemes
// @Description Update authentication settings, this api is available only to administrators
// @Tags Templates
// @Accept json
// @Produce json
// @Param request body HandleUpdateServerSettingsRequestBody true "authentication settings"
// @Success 200 {object} serializers.AuthenticationSettingsSerializer
// @Failure 400 "Bad request"
// @Failure 406 "Email server is not configured"
// @Failure 500 "Internal server error"
// @Router /api/v1/admin/authentication-settings [put]
func HandleUpdateAuthenticationSettings(c *gin.Context) {
	if !config.IsEmailConfigured() {
		utils.ErrorResponse(
			c,
			http.StatusNotAcceptable,
			"email server is not configured, instance settings cannot be updated",
		)
		return
	}

	// parse and validate request body
	var parsedBody HandleUpdateServerSettingsRequestBody
	err := c.ShouldBindBodyWithJSON(&parsedBody)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing or invalid field")
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

	s.IsSignUpOpen = *parsedBody.IsSignUpOpen
	s.IsSignUpRestricted = *parsedBody.IsSignUpRestricted
	s.AllowedEmailRegex = *parsedBody.AllowedEmailRegex
	s.BlockedEmailRegex = *parsedBody.BlockedEmailRegex
	s.UsersMustBeApproved = *parsedBody.UsersMustBeApproved
	s.ApprovedByDefaultEmailRegex = *parsedBody.ApprovedByDefaultEmailRegex
	models.SaveSingletonModel(s)

	c.JSON(http.StatusOK, serializers.LoadAuthenticationSettingsSerializer(s))
}
