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

func generateVerificationCode(user User) string {
	hasher := sha1.New()
	hasher.Write([]byte(strconv.Itoa(int(time.Now().UnixNano())) + user.Email))
	code := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return code
}

/*
CreteEmailVerificationCode creates a new email verification code for
the given user with an optional expiration time.
*/
func CreateEmailVerificationCode(
	expiration *time.Time,
	user User,
) (*EmailVerificationCode, error) {
	// generate a unique verification code
	code := generateVerificationCode(user)
	vc, err := RetrieveVerificationCodeByCode(code)
	if err != nil {
		return nil, err
	}

	for vc != nil {
		code := generateVerificationCode(user)
		vc, err = RetrieveVerificationCodeByCode(code)
		if err != nil {
			return nil, err
		}
	}

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

/*
RevokeAllTokensForUser revokes all email verification
codes for the given user.
*/
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

/*
RetrieveVerificationCodeByCode retrieves an email verification code
by its code.
*/
func RetrieveVerificationCodeByCode(code string) (*EmailVerificationCode, error) {
	vc := EmailVerificationCode{}
	r := dbconn.DB.Preload("User").First(
		&vc,
		map[string]interface{}{
			"code": code,
		},
	)

	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, r.Error
	}

	return &vc, nil
}
