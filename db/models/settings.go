package models

import (
	"errors"
	"time"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

type InstanceSettings struct {
	ID                 uint           `gorm:"primarykey" json:"id"`
	IsSignUpOpen       bool           `gorm:"column:is_signup_open; default:false"`
	IsSignUpRestricted bool           `gorm:"column:is_signup_restricted; default:false"`
	AllowedEmailRegex  string         `gorm:"column:allowed_email_regex; type:text;"`
	BlockedEmailRegex  string         `gorm:"column:blocked_email_regex; type:text;"`
	CreatedAt          time.Time      `json:"-"`
	UpdatedAt          time.Time      `json:"-"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

/*
 */
func (s *InstanceSettings) UpdateInstanceSettings() error {
	s.ID = 1

	if err := dbconn.DB.
		Save(s).Error; err != nil {
		return err
	}

	return nil
}

/*
Retrieve instance settings, if not found return default value
*/
func GetInstanceSettings() (*InstanceSettings, error) {
	var settings InstanceSettings
	if err := dbconn.DB.First(&settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &settings, nil
		}
		return nil, err
	}

	return &settings, nil
}
