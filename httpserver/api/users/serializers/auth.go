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

type IsSignUpOpenSerializer struct {
	IsSignUpOpen bool `json:"is_signup_open"`
}

func LoadIsSignUpOpenSerializer(isSignUpOpen bool) IsSignUpOpenSerializer {
	return IsSignUpOpenSerializer{
		IsSignUpOpen: isSignUpOpen,
	}
}

type CanResetPasswordSerializer struct {
	CanResetPassword bool `json:"can_reset_password"`
}

func LoadCanResetPasswordSerializer(canResetPassword bool) CanResetPasswordSerializer {
	return CanResetPasswordSerializer{
		CanResetPassword: canResetPassword,
	}
}

type RequestPasswordResetSerializer struct {
	Success bool   `json:"success"`
	Email   string `json:"email"`
}

func LoadRequestPasswordResetSerializer(success bool, email string) RequestPasswordResetSerializer {
	return RequestPasswordResetSerializer{
		Success: success,
		Email:   email,
	}
}
