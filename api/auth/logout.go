package auth

import (
	"net/http"
	"strings"

	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

func HandleLogout(ctx *gin.Context) {

	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": "missing or invalid authorization token",
		})
	}

	headerParts := strings.Split(authHeader, "Bearer ")

	if len(headerParts) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": "missing or invalid authorization token",
		})
	}

	jwtToken := headerParts[1]

	var token models.Token
	result := db.DB.Where("token=?", jwtToken).Preload("User").First(&token)
	if result.Error != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": "missing or invalid authorization token",
		})
	}

	db.DB.Delete(&token)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
