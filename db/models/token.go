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
	ID                 uint       `gorm:"primarykey"`
	Token              string     `gorm:"column:token; size:255;unique;"`
	ExpirationDate     *time.Time `gorm:"column:expiration_date;"`
	UserID             uint       `gorm:"column:user_id;"`
	User               User       `gorm:"constraint:OnDelete:CASCADE;"`
	ImpersonatedUserID uint       `gorm:"column:impersonated_user_id;"`
	ImpersonatedUser   *User      `gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
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

/*
GetLoginCountPerDayInLast7Days returns an array of login counts for each of the last 7 days.
The array is ordered from oldest to newest day.
*/
func GetLoginCountPerDayInLast7Days() ([]int64, error) {
	var counts []int64
	today := time.Now().Truncate(24 * time.Hour)

	for i := 6; i >= 0; i-- {
		day := today.AddDate(0, 0, -i)
		nextDay := day.AddDate(0, 0, 1)

		var count int64
		err := dbconn.DB.Model(&Token{}).
			Where("created_at >= ? AND created_at < ?", day, nextDay).
			Count(&count).Error
		if err != nil {
			return nil, err
		}

		counts = append(counts, count)
	}

	return counts, nil
}

/*
GetLastLoginTimeForUser retrieves the last login time for a specific user.
If the user has never logged in, it returns nil without an error.
*/
func GetLastLoginTimeForUser(user *User) (*time.Time, error) {
	var token Token
	err := dbconn.DB.Where("user_id = ?", user.ID).
		Order("created_at DESC").
		First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No login records found
		}
		return nil, err
	}

	return &token.CreatedAt, nil
}
