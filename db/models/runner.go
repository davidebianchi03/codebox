package models

import (
	"time"

	"gorm.io/gorm"
)

type Runner struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Name          string         `gorm:"size:255;unique;not null;" json:"name"`
	Token         string         `gorm:"size:255;unique;not null;" json:"-"`
	Port          uint           `gorm:"default:0;" json:"-"`
	Type          string         `gorm:"size:255;" json:"type"`
	Restricted    bool           `gorm:"default:false;" json:"-"`
	AllowedGroups []Group        `gorm:"many2many:runner_allowed_groups;" json:"-"`
	UsePublicUrl  bool           `gorm:"default:false;" json:"use_public_url"`
	PublicUrl     string         `gorm:"type:text;" json:"public_url"`
	LastContact   *time.Time     `json:"last_contact"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
