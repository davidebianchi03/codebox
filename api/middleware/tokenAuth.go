package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"codebox.com/db"
	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware(ctx *gin.Context) {
	user, err := GetUserFromContext(ctx)
	_ = user
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"detail": err.Error(),
		})
		return
	}
	ctx.Next()
}

func GetUserFromContext(ctx *gin.Context) (db.User, error) {
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		return db.User{}, fmt.Errorf("missing or invalid authorization token")
	}

	headerParts := strings.Split(authHeader, "Bearer ")

	if len(headerParts) != 2 {
		return db.User{}, fmt.Errorf("missing or invalid authorization token")
	}

	jwtToken := headerParts[1]

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
