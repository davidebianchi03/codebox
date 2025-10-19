package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
)

// retrieves token from 'Authorization' header
func getTokenFromAuthorizationHeader(ctx *gin.Context) (token string, err error) {
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader != "" {
		headerParts := strings.Split(authHeader, "Bearer ")

		if len(headerParts) != 2 {
			return "", fmt.Errorf("missing or invalid authorization token")
		}

		return headerParts[1], nil
	}
	return "", nil
}

// retrieves token from authorization cookie
func getTokenFromAuthCookie(ctx *gin.Context) (token string, err error) {
	cookie, err := ctx.Cookie(config.Environment.AuthCookieName)
	if err != nil {
		return "", err
	}
	return cookie, nil
}

// Retrieves the authorization token from the cookie used in subdomains.
// Note: This cookie uses a different name compared to the main website's authorization cookie.
// This addresses a scenario where a subdomain (within the codebox server's wildcard domain)
// might attempt to set a cookie with the same name as the secure codebox server's
// authorization cookie, which browsers prevent.
func getTokenFromSubdomainAuthCookie(ctx *gin.Context) (token string, err error) {
	cookie, err := ctx.Cookie(config.Environment.SubdomainAuthCookieName)
	if err != nil {
		return "", err
	}
	return cookie, nil
}

// retrieve authorization token from the gin context
// sources hierarchy is:
// 1. authorization header
// 2. main website's authorization cookie
// 3. subdomains' authorization cookie
func GetTokenFromContext(ctx *gin.Context) (models.Token, error) {
	t, _ := getTokenFromAuthorizationHeader(ctx)
	if t == "" {
		t, _ = getTokenFromAuthCookie(ctx)
		if t == "" {
			t, _ = getTokenFromSubdomainAuthCookie(ctx)
		}
	}

	if t == "" {
		return models.Token{}, fmt.Errorf("missing or invalid authorization token")
	}

	var token models.Token
	result := dbconn.DB.Where("token=?", t).Preload("User").Preload("ImpersonatedUser").First(&token)
	if result.Error != nil {
		return models.Token{}, fmt.Errorf("missing or invalid authorization token")
	}

	if token.ExpirationDate.Compare(time.Now()) == -1 {
		return models.Token{}, fmt.Errorf("missing or invalid authorization token")
	}

	return token, nil
}
