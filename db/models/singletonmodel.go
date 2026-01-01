package models

import (
	"errors"
	"time"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

/*
Singleton model, a db table that has only one record
*/
type SingletonModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (s *SingletonModel) SetID(id uint) {
	s.ID = id
}

func (s *SingletonModel) GetID() uint {
	return s.ID
}

/*
Get the instance of a singleton model
*/
func GetSingletonModelInstance[T any]() (*T, error) {
	var m T
	if err := dbconn.DB.First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &m, nil
		}
		return nil, err
	}

	return &m, nil
}

type HasID interface {
	SetID(uint)
	GetID() uint
}

/*
Save a singleton model, instances of a singleton model
always have as id 1
*/
func SaveSingletonModel[T HasID](s T) error {
	s.SetID(1)
	if err := dbconn.DB.
		Save(s).Error; err != nil {
		return err
	}

	return nil
}
