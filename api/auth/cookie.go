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

// set authentication cookie for subdomains
// Note: This cookie uses a different name compared to the main website's authorization cookie.
// This addresses a scenario where a subdomain (within the codebox server's wildcard domain)
// might attempt to set a cookie with the same name as the secure codebox server's
// authorization cookie, which browsers prevent.
func SetSubdomainsAuthCookie(ctx *gin.Context, token string) error {
	// Set auth cookie, duration is set to zero because
	// token expiration has been already set in DB
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		config.Environment.SubdomainAuthCookieName,
		token,
		0,
		"",
		"",
		false,
		false,
	)

	return nil
}
