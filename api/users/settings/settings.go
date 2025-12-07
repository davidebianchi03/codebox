package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/db/models"
)

// HandleRetrieveServerSettings godoc
// @Summary Retrieve instance settings
// @Schemes
// @Description Retrieve instance settings, this api is available only to administrators
// @Tags Templates
// @Accept json
// @Produce json
// @Success 200 {object} serializers.InstanceSettingsSerializer
// @Router /api/v1/admin/instance-settings [get]
func HandleRetrieveServerSettings(c *gin.Context) {
	is, err := models.GetInstanceSettings()
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	c.JSON(http.StatusOK, serializers.LoadInstanceSettingsSerializer(is))
}

type HandleUpdateServerSettingsRequestBody struct {
	IsSignUpOpen       bool   `json:"is_signup_open"`
	IsSignUpRestricted bool   `json:"is_signup_restricted"`
	AllowedEmailRegex  string `json:"allowed_emails_regex"`
	BlockedEmailRegex  string `json:"blocked_emails_regex"`
}

// HandleUpdateServerSettings godoc
// @Summary Update instance settings
// @Schemes
// @Description Update instance settings, this api is available only to administrators
// @Tags Templates
// @Accept json
// @Produce json
// @Param request body HandleUpdateServerSettingsRequestBody true "Instance settings"
// @Success 200 {object} serializers.InstanceSettingsSerializer
// @Router /api/v1/admin/instance-settings [put]
func HandleUpdateServerSettings(c *gin.Context) {
	// parse and validate request body
	var parsedBody HandleUpdateServerSettingsRequestBody
	err := c.ShouldBindBodyWithJSON(&parsedBody)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	is, err := models.GetInstanceSettings()
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	is.IsSignUpOpen = parsedBody.IsSignUpOpen
	is.IsSignUpRestricted = parsedBody.IsSignUpRestricted
	is.AllowedEmailRegex = parsedBody.AllowedEmailRegex
	is.BlockedEmailRegex = parsedBody.BlockedEmailRegex
	is.UpdateInstanceSettings()

	c.JSON(http.StatusOK, serializers.LoadInstanceSettingsSerializer(is))
}
