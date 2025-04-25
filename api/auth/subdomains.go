package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/davidebianchi03/codebox/api/utils"
	"github.com/davidebianchi03/codebox/config"
	dbconn "github.com/davidebianchi03/codebox/db/connection"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

func HandleSubdomainLoginAuthorize(ctx *gin.Context) {
	next, ok := ctx.GetQuery("next")
	if !ok {
		ctx.Redirect(
			http.StatusTemporaryRedirect,
			"",
		)
		return
	}

	nextUrl, err := url.Parse(next)
	if err != nil {
		ctx.Redirect(
			http.StatusTemporaryRedirect,
			"",
		)
		return
	}

	token, err := utils.GetTokenFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	authorizationCode, err := models.GenerateAuthorizationCode(token, time.Now().Add(2*time.Minute))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
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

	ctx.Redirect(
		http.StatusTemporaryRedirect,
		redirectUri,
	)
}

func HandleSubdomainLoginCallback(ctx *gin.Context) {
	code, ok := ctx.GetQuery("code")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing or invalid authorization code",
		})
		return
	}

	next, ok := ctx.GetQuery("next")
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing or invalid authorization code",
		})
		return
	}

	// check if token is expired
	// if token is expired delete it and return an error message
	if authorizationCode.ExpiresAt.UnixMilli() < time.Now().UnixMilli() {
		dbconn.DB.Unscoped().Delete(&authorizationCode)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing or invalid authorization code",
		})
		return
	}

	// set cookie
	if err := SetAuthCookie(ctx, authorizationCode.Token.Token); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	// remove expired authorization codes from db
	models.RemoveExpiredAuthorizationCodes()

	ctx.Redirect(
		http.StatusTemporaryRedirect,
		next,
	)
}
