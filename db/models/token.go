package models

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

var secretKey = []byte("secret-key") // TODO: replace on build

type Token struct {
	gorm.Model
	ID             uint       `gorm:"primarykey"`
	Token          string     `gorm:"column:token; size:255;unique;"`
	ExpirationDate *time.Time `gorm:"column:expiration_date;"`
	UserID         uint       `gorm:"column:user_id;"`
	User           User       `gorm:"constraint:OnDelete:CASCADE;"`
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

func CreateToken(user User, duration time.Duration) (Token, error) {
	tokenExpiration := time.Now().Add(duration)

	jwtToken, err := generateJWTToken(user.ID, tokenExpiration)

	if err != nil {
		return Token{}, fmt.Errorf("cannot create token, %s", err)
	}

	token := Token{
		Token:          jwtToken,
		ExpirationDate: &tokenExpiration,
		User:           user,
	}

	if err := dbconn.DB.Create(&token).Error; err != nil {
		return Token{}, err
	}

	return token, nil
}
