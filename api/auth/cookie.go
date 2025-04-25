package auth

import (
	"net/http"

	"github.com/davidebianchi03/codebox/config"
	"github.com/gin-gonic/gin"
)

// set authentication cookie
func SetAuthCookie(ctx *gin.Context, token string) error {
	// Set auth cookie, duration is set to zero because
	// token expiration has been already set in DB
	ctx.SetSameSite(http.SameSiteLaxMode) // TODO: use same site lax in prod
	ctx.SetCookie(
		config.Environment.AuthCookieName,
		token,
		0,
		"",
		"",
		false,
		false,
	)

	return nil
}
