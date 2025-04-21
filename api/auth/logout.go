package auth

import (
	"fmt"
	"net/http"

	"github.com/davidebianchi03/codebox/api/utils"
	"github.com/davidebianchi03/codebox/db"
	"github.com/gin-gonic/gin"
)

func HandleLogout(ctx *gin.Context) {
	token, err := utils.GetTokenFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"detail": err.Error(),
		})
	}

	db.DB.Delete(&token)

	// clear cookies
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(
		"jwtToken",
		"",
		3600*24*20,
		"",
		ctx.Request.Host,
		true,
		false,
	)
	ctx.SetCookie(
		"jwtToken",
		"",
		3600*24*20,
		"",
		fmt.Sprintf(".%s", ctx.Request.Host),
		true,
		false,
	)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
