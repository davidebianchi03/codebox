package email

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/emails"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
)

// HandleSendTestEmail godoc
// @Summary Send Test Email
// @Schemes
// @Description Send an email to test email send service, the email is sent synchronously so it can take few seconds
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} serializers.TestEmailResponseSerializer
// @Router /api/v1/admin/send-test-email [post]
func HandleSendTestEmail(c *gin.Context) {
	currentUser, err := utils.GetUserFromContext(c)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	if err := emails.SendTestEmail(
		currentUser.Email,
	); err != nil {
		c.JSON(
			http.StatusOK,
			serializers.LoadTestEmailResponseSerializer(
				false, err.Error(),
			),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		serializers.LoadTestEmailResponseSerializer(
			true, "",
		),
	)
}
