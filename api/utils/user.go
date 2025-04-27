package utils

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/db/models"
)

// get the current user from the request context
func GetUserFromContext(ctx *gin.Context) (models.User, error) {
	token, err := GetTokenFromContext(ctx)
	if err != nil {
		return models.User{}, err
	}

	return token.User, nil
}
