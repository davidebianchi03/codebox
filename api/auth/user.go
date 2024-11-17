package auth

import (
	"net/http"

	"codebox.com/api/utils"
	"codebox.com/db"
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

func HandleUpdateUserDetails(ctx *gin.Context) {

	type RequestBody struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	var parsedBody RequestBody
	err := ctx.ShouldBindBodyWithJSON(&parsedBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": err.Error(),
		})
		return
	}

	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if parsedBody.FirstName != "" {
		user.FirstName = parsedBody.FirstName
	}
	if parsedBody.LastName != "" {
		user.LastName = parsedBody.LastName
	}

	db.DB.Save(&user)

	ctx.JSON(http.StatusOK, gin.H{
		"email":        user.Email,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"is_superuser": user.IsSuperuser,
		"public_key":   user.SshPublicKey,
	})
}
