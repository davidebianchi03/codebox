package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/db/models"
)

/*
Method that returns if the user is being impersonated
*/
func UserIsBeingImpersonated(c *gin.Context) (bool, error) {
	t, err := utils.GetTokenFromContext(c)
	if err != nil {
		return false, err
	}

	u, err := utils.GetUserFromContext(c)
	if err != nil {
		return false, err
	}

	return t.User.Email != u.Email, nil
}

// HandleStopImpersonation godoc
// @Summary API to stop the impersonation of a user
// @Schemes
// @Description API to stop the impersonation of a user
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200
// @Router /api/v1/stop-impersonation [post]
func HandleStopImpersonation(c *gin.Context) {
	token, err := utils.GetTokenFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if token.ImpersonatedUser == nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"no user is being impersonated in this session",
		)
		return
	}

	log, err := models.RetrieveLatestImpersonationLogByToken(token)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	if log != nil {
		now := time.Now()
		log.ImpersonationFinishedAt = &now
		if err := models.UpdateImpersonationLog(log); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	token.ImpersonatedUserID = nil
	token.ImpersonatedUser = nil
	if err := models.UpdateToken(token); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"detail": "impersonation has been stopped",
	})
}
