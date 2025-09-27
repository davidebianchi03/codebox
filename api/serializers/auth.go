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
