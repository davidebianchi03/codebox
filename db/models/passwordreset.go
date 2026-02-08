package models

import (
	"crypto/md5"
	"fmt"
	"time"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

/*
Model for password reset tokens
*/
type PasswordResetToken struct {
	ID         uint      `gorm:"primarykey" json:"-"`
	UserID     uint      `gorm:"column:user_id; not null; index;" json:"-"`
	User       User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Token      string    `gorm:"column:token; size:255; unique; not null;" json:"-"`
	Expiration time.Time `gorm:"column:expiration; not null;" json:"-"`
	CreatedAt  int64     `gorm:"column:created_at; autoCreateTime" json:"-"`
}

/*
Generate a password reset token for a given user
*/
func CreatePasswordResetToken(user User) (*PasswordResetToken, error) {
	expiration := time.Now().Add(24 * time.Hour)
	plainTokenContent := fmt.Sprintf("%s-%d-%d", user.Email, user.ID, time.Now().Unix())

	// Generate a hash of the plain token content
	token := fmt.Sprintf("%x", md5.Sum([]byte(plainTokenContent)))

	prt := PasswordResetToken{
		UserID:     user.ID,
		Token:      token,
		Expiration: expiration,
	}

	r := dbconn.DB.Create(&prt)
	if r.Error != nil {
		return nil, r.Error
	}

	return &prt, nil
}

/*
Generate a password reset token for a given user
*/
func GetPasswordResetToken(token string) (*PasswordResetToken, error) {
	var prt PasswordResetToken

	r := dbconn.DB.
		Preload("User").
		Where("token = ?", token).
		First(&prt)

	if r.Error != nil {
		if r.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, r.Error
	}

	return &prt, nil
}

/*
Delete a password reset token
*/
func DeletePasswordResetToken(prt PasswordResetToken) error {
	r := dbconn.DB.Delete(&prt)
	return r.Error
}

/*
DeleteExpiredPasswordResetTokens deletes all expired password
reset tokens from the database
*/
func DeleteExpiredPasswordResetTokens() error {
	r := dbconn.DB.
		Where("expiration < ?", time.Now()).
		Delete(&PasswordResetToken{})
	return r.Error
}
