package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"

	httperrors "gitlab.com/codebox4073715/codebox/httpserver/errors"
)

func HandleSubdomainLoginAuthorize(c *gin.Context) {
	next, ok := c.GetQuery("next")
	if !ok {
		c.Redirect(
			http.StatusTemporaryRedirect,
			"",
		)
		return
	}

	nextUrl, err := url.Parse(next)
	if err != nil {
		c.Redirect(
			http.StatusTemporaryRedirect,
			"",
		)
		return
	}

	token, err := utils.GetTokenFromContext(c)
	if err != nil {
		httperrors.RenderError(
			c,
			http.StatusInternalServerError,
			"Internal server error",
		)
		return
	}

	authorizationCode, err := models.GenerateAuthorizationCode(token, time.Now().Add(2*time.Minute))
	if err != nil {
		httperrors.RenderError(
			c,
			http.StatusInternalServerError,
			"Internal server error",
		)
		return
	}

	redirectUri := fmt.Sprintf(
		"%s://%s/api/v1/auth/subdomains/callback-%s?code=%s&next=%s",
		nextUrl.Scheme,
		nextUrl.Host,
		url.PathEscape(config.Environment.AuthCookieName),
		url.QueryEscape(authorizationCode.Code),
		url.QueryEscape(next),
	)

	c.Redirect(
		http.StatusTemporaryRedirect,
		redirectUri,
	)
}

func HandleSubdomainLoginCallback(c *gin.Context) {
	code, ok := c.GetQuery("code")
	if !ok {
		httperrors.RenderError(
			c, http.StatusBadRequest, "missing or invalid authorization code",
		)
		return
	}

	next, ok := c.GetQuery("next")
	if !ok {
		next = ""
	}

	// retrive authorization code from db
	var authorizationCode *models.AuthorizationCode
	if err := dbconn.DB.Preload("Token").Find(
		&authorizationCode,
		map[string]interface{}{
			"code": code,
		},
	).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing or invalid authorization code",
		})
		return
	}

	// check if token is expired
	// if token is expired delete it and return an error message
	if authorizationCode.ExpiresAt.UnixMilli() < time.Now().UnixMilli() {
		dbconn.DB.Unscoped().Delete(&authorizationCode)
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing or invalid authorization code",
		})
		return
	}

	// set cookie
	if err := SetSubdomainsAuthCookie(c, authorizationCode.Token.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	// remove expired authorization codes from db
	models.RemoveExpiredAuthorizationCodes()

	c.Redirect(
		http.StatusTemporaryRedirect,
		next,
	)
}
