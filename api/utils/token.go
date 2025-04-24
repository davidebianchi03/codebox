package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/davidebianchi03/codebox/config"
	dbconn "github.com/davidebianchi03/codebox/db/connection"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

func GetTokenFromContext(ctx *gin.Context) (models.Token, error) {
	authHeader := ctx.Request.Header.Get("Authorization")

	jwtCookie, err := ctx.Cookie(config.Environment.AuthCookieName)
	if err != nil {
		jwtCookie = ""
	}

	if authHeader == "" && jwtCookie == "" {
		return models.Token{}, fmt.Errorf("missing or invalid authorization token")
	}

	jwtToken := ""
	if authHeader != "" {
		headerParts := strings.Split(authHeader, "Bearer ")

		if len(headerParts) != 2 {
			return models.Token{}, fmt.Errorf("missing or invalid authorization token")
		}

		jwtToken = headerParts[1]
	} else {
		jwtToken = jwtCookie
	}

	var token models.Token
	result := dbconn.DB.Where("token=?", jwtToken).Preload("User").First(&token)
	if result.Error != nil {
		return models.Token{}, fmt.Errorf("missing or invalid authorization token")
	}

	if token.ExpirationDate.Compare(time.Now()) == -1 {
		return models.Token{}, fmt.Errorf("missing or invalid authorization token")
	}

	return token, nil
}
