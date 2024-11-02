package auth

import (
	"net/http"

	"codebox.com/api/utils"
	"github.com/gin-gonic/gin"
)

func HandleRetriveUserDetails(ctx *gin.Context) {
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"email":        user.Email,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"is_superuser": user.IsSuperuser,
		"public_key":   user.SshPublicKey,
	})
}
