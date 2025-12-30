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

/*
Save a singleton model, instances of a singleton model
always have as id 1
*/
func (s *SingletonModel) SaveSingletonModel() error {
	s.ID = 1

	if err := dbconn.DB.
		Save(s).Error; err != nil {
		return err
	}

	return nil
}
