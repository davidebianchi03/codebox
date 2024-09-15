package db

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var secretKey = []byte("secret-key") // TODO: replace on build

type Token struct {
	gorm.Model
	Token          string    `gorm:"column:token; size:1024;unique;"`
	ExpirationDate time.Time `gorm:"column:expiration_date;"`
	UserID         uint      `gorm:"column:user_id;"`
	User           User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
		return Token{}, fmt.Errorf("Cannot create token, %s", err)
	}

	token := Token{
		Token:          jwtToken,
		ExpirationDate: tokenExpiration,
		User:           user,
	}
	return token, nil
}
