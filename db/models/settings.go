package models

import (
	"time"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

type InstanceSettings struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	AllowUserSignUp   bool           `gorm:"column:allow_user_sign_up; default:false"`
	SignUpRestricted  bool           `gorm:"column:sign_up_restricted; default:false"`
	AllowedEmailRegex string         `gorm:"column:allowed_email_regex; type:text;"`
	BlockedEmailRegex string         `gorm:"column:blocked_email_regex; type:text;"`
	CreatedAt         time.Time      `json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
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
 */
func GetInstanceSettings() (*InstanceSettings, error) {
	var settings InstanceSettings
	if err := dbconn.DB.
		First(&settings).Error; err != nil {
		return nil, err
	}

	return &settings, nil
}
