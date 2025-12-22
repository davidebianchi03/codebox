package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/db/models"
)

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
