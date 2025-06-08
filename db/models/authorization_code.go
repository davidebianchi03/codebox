package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

type AuthorizationCode struct {
	ID        uint   `gorm:"primarykey"`
	Code      string `gorm:"column:code; size:255"`
	TokenID   uint   `gorm:"column:token_id;"`
	Token     *Token
	ExpiresAt time.Time      `gorm:"column:expires_at;"`
	CreatedAt time.Time      `gorm:"column:created_at;"`
	UpdatedAt time.Time      `gorm:"column:updated_at;"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// generate authorization code
func GenerateAuthorizationCode(token Token, expiration time.Time) (ac AuthorizationCode, err error) {
	// generate a random code
	// if it try to generate a new one
	code := ""
	exists := true

	for exists {
		b := make([]byte, 32)
		rand.Read(b)
		code = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)

		r := dbconn.DB.Model(AuthorizationCode{}).Where("code = ?", code)
		if r.Error != nil {
			return AuthorizationCode{}, errors.New("unable to check if token alredy exists")
		}

		exists = r.RowsAffected > 0
	}

	// store authorization token in db
	ac = AuthorizationCode{
		Code:      code,
		TokenID:   token.ID,
		Token:     &token,
		ExpiresAt: expiration,
	}

	if err := dbconn.DB.Save(&ac).Error; err != nil {
		return AuthorizationCode{}, errors.New("unable to create token")
	}
	return ac, nil
}

// remove expired authorization codes
func RemoveExpiredAuthorizationCodes() error {
	// list expired authorization codes
	expiredAuthorizationCodes := []AuthorizationCode{}
	now := time.Now()
	if err := dbconn.DB.Where("expires_at < ?", now).Find(&expiredAuthorizationCodes).Error; err != nil {
		return errors.New("failed to list expired authorization codes")
	}

	// delete codes
	for _, c := range expiredAuthorizationCodes {
		dbconn.DB.Unscoped().Delete(&c)
	}

	return nil
}
