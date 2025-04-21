package utils

import (
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/gin-gonic/gin"
)

func GetUserFromContext(ctx *gin.Context) (models.User, error) {
	token, err := GetTokenFromContext(ctx)
	if err != nil {
		return models.User{}, err
	}

	return token.User, nil
}
