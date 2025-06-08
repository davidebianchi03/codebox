package models

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID        uint           `gorm:"primarykey"`
	Name      string         `gorm:"column:name; size:255;unique"`
	CreatedAt time.Time      `gorm:"column:created_at;"`
	UpdatedAt time.Time      `gorm:"column:updated_at;"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
