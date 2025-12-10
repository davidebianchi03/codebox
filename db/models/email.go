package models

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"strconv"
	"time"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

type EmailVerificationCode struct {
	ID         uint           `gorm:"primarykey" json:"-"`
	Code       string         `gorm:"column:code; size:255; unique; not null;"`
	Expiration *time.Time     `gorm:"column:expiration;"`
	UserID     uint           `gorm:"column:user_id;"`
	User       User           `gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func CreateEmailVerificationCode(
	expiration *time.Time,
	user User,
) (*EmailVerificationCode, error) {
	hasher := sha1.New()
	hasher.Write([]byte(strconv.Itoa(int(time.Now().UnixNano())) + user.Email))
	code := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	emailVerificationCode := &EmailVerificationCode{
		Code:       code,
		Expiration: expiration,
		UserID:     user.ID,
		User:       user,
	}

	result := dbconn.DB.Create(emailVerificationCode)
	if result.Error != nil {
		return nil, result.Error
	}

	return emailVerificationCode, nil
}

func RevokeAllTokensForUser(user User) error {
	r := dbconn.DB.Unscoped().Delete(
		EmailVerificationCode{},
		map[string]interface{}{
			"user_id": user.ID,
		},
	)

	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		return r.Error
	}

	return nil
}
