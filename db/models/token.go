package models

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var secretKey = []byte("secret-key") // TODO: replace on build

type Token struct {
	gorm.Model
	ID             uint   `gorm:"primarykey"`
	Token          string `gorm:"size:1024;unique;"`
	ExpirationDate time.Time
	UserID         uint
	User           User `gorm:"constraint:OnDelete:CASCADE;"`
}

func generateJWTToken(userId uint, expiration time.Time) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  userId,
			"exp": expiration,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateToken(user User) (Token, error) {
	tokenExpiration := time.Now().Add(time.Duration(time.Hour * 24))

	jwtToken, err := generateJWTToken(user.ID, tokenExpiration)

	if err != nil {
		return Token{}, fmt.Errorf("cannot create token, %s", err)
	}

	token := Token{
		Token:          jwtToken,
		ExpirationDate: tokenExpiration,
		User:           user,
	}
	return token, nil
}
