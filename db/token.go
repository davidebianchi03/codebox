package db

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	Id             uint   `gorm:"unique;primaryKey;autoIncrement"`
	Token          string `gorm:"size:1024;unique;"`
	ExpirationDate time.Time
	User           User `gorm:"foreignKey:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
