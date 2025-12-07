package models

import (
	"time"

	"gorm.io/gorm"
)

type EmailVerificationCode struct {
	ID         uint           `gorm:"primarykey" json:"-"`
	Code       string         `gorm:"column:code; size:255; unique; not null;"`
	Expiration time.Time      `gorm:"column:expiration;"`
	UserID     uint           `gorm:"column:user_id;"`
	User       User           `gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
