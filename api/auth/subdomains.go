package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	apierrors "gitlab.com/codebox4073715/codebox/api/errors"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
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
		apierrors.RenderError(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	authorizationCode, err := models.GenerateAuthorizationCode(token, time.Now().Add(2*time.Minute))
	if err != nil {
		apierrors.RenderError(ctx, http.StatusInternalServerError, "Internal server error")
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
		apierrors.RenderError(
			ctx, http.StatusBadRequest, "missing or invalid authorization code",
		)
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
	if err := SetSubdomainsAuthCookie(ctx, authorizationCode.Token.Token); err != nil {
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
