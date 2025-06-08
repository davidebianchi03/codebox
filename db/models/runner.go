package models

import (
	"time"

	"gorm.io/gorm"
)

type Runner struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	Name          string         `gorm:"column:name; size:255;unique;not null;" json:"name"`
	Token         string         `gorm:"column:token; size:255;unique;not null;" json:"-"`
	Port          uint           `gorm:"column:port; default:0;" json:"-"`
	Type          string         `gorm:"column:type; size:255;" json:"type"`
	Restricted    bool           `gorm:"column:restricted; default:false;" json:"-"`
	AllowedGroups []Group        `gorm:"many2many:runner_allowed_groups;" json:"-"`
	UsePublicUrl  bool           `gorm:"column:use_public_url; default:false;" json:"use_public_url"`
	PublicUrl     string         `gorm:"column:public_url; type:text;" json:"public_url"`
	LastContact   *time.Time     `gorm:"column:last_contact;" json:"last_contact"`
	CreatedAt     time.Time      `gorm:"column:created_at;" json:"-"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;" json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
