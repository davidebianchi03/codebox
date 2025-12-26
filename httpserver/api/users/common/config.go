package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/serializers"
)

// HandleAdminEmailServiceConfigured godoc
// @Summary Check if email service is configured
// @Schemes
// @Description Check if email service is configured, this api is available only to administrators
// @Tags Templates
// @Accept json
// @Produce json
// @Success 200 {object} serializers.EmailServiceConfiguredSerializer
// @Router /api/v1/admin/email-service-configured [get]
func HandleAdminEmailServiceConfigured(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		serializers.LoadEmailServiceConfiguredSerializer(config.IsEmailConfigured()),
	)
}
