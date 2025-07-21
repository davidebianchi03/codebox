package serializers

import (
	"time"

	"gitlab.com/codebox4073715/codebox/db/models"
)

type TokenSerializer struct {
	Token          string    `json:"token"`
	ExpirationDate time.Time `json:"expiration"`
}

func LoadTokenSerializer(token *models.Token) *TokenSerializer {
	if token == nil {
		return nil
	}

	return &TokenSerializer{
		Token:          token.Token,
		ExpirationDate: *token.ExpirationDate,
	}
}

func LoadMultipleTokenSerializer(tokens []models.Token) []TokenSerializer {
	serializers := make([]TokenSerializer, len(tokens))
	for i, token := range tokens {
		serializers[i] = *LoadTokenSerializer(&token)
	}
	return serializers
}

type UserSerializer struct {
	Email             string `json:"email"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	IsSuperUser       bool   `json:"is_superuser"`
	IsTemplateManager bool   `json:"is_template_manager"`
}

func LoadUserSerializer(user *models.User) *UserSerializer {
	return &UserSerializer{
		Email:             user.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		IsSuperUser:       user.IsSuperuser,
		IsTemplateManager: user.IsTemplateManager,
	}
}

func LoadMultipleUserSerializer(users []models.User) []UserSerializer {
	serializers := make([]UserSerializer, len(users))
	for i, user := range users {
		serializers[i] = *LoadUserSerializer(&user)
	}
	return serializers
}
