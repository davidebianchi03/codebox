package utils

import (
	"fmt"
	"strings"
	"time"

	"codebox.com/db"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(ctx *gin.Context) (db.User, error) {
	authHeader := ctx.Request.Header.Get("Authorization")

	jwtCookie, err := ctx.Cookie("jwtToken")
	if err != nil {
		jwtCookie = ""
	}

	if authHeader == "" && jwtCookie == "" {
		return db.User{}, fmt.Errorf("missing or invalid authorization token")
	}

	jwtToken := ""
	if authHeader != "" {
		headerParts := strings.Split(authHeader, "Bearer ")

		if len(headerParts) != 2 {
			return db.User{}, fmt.Errorf("missing or invalid authorization token")
		}

		jwtToken = headerParts[1]
	} else {
		jwtToken = jwtCookie
	}

	var token db.Token
	result := db.DB.Where("token=?", jwtToken).Preload("User").First(&token)
	if result.Error != nil {
		return db.User{}, fmt.Errorf("missing or invalid authorization token")
	}

	if token.ExpirationDate.Compare(time.Now()) == -1 {
		return db.User{}, fmt.Errorf("missing or invalid authorization token")
	}

	return token.User, nil
}
