package models

import (
	"time"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

type ImpersonationLog struct {
	ID                      uint       `gorm:"primarykey"`
	TokenID                 *uint      `gorm:"column:token_id;"`
	Token                   *Token     `gorm:"constraint:OnDelete:SET NULL;"`
	ImpersonatorID          uint       `gorm:"column:impersonated_user_id;not null;"`
	Impersonator            User       `gorm:"constraint:OnDelete:CASCADE;"`
	ImpersonatorIPAddress   string     `gorm:"column:impersonator_ip_address;not null;"`
	ImpersonatedUserID      uint       `gorm:"column:impersonated_user_id;not null;"`
	ImpersonatedUser        User       `gorm:"constraint:OnDelete:CASCADE;"`
	ImpersonationStartedAt  time.Time  `gorm:"column:impersonation_started_at;not null;"`
	ImpersonationFinishedAt *time.Time `gorm:"column:impersonation_finished_at;"`
	UpdatedAt               time.Time
	DeletedAt               gorm.DeletedAt `gorm:"index"`
}

/*
CreateImpersonationLog create an impersonation log
*/
func CreateImpersonationLog(
	token Token,
	impersonator User,
	impersonatorIP string,
	impersonatedUser User,
) (*ImpersonationLog, error) {
	l := ImpersonationLog{
		Token:                  &token,
		Impersonator:           impersonator,
		ImpersonatorIPAddress:  impersonatorIP,
		ImpersonatedUser:       impersonatedUser,
		ImpersonationStartedAt: time.Now(),
	}

	r := dbconn.DB.Create(&l)
	if r.Error != nil {
		return nil, r.Error
	}

	return &l, nil
}
