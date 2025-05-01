package models

import (
	"os"
	"path"
	"time"

	"gitlab.com/codebox4073715/codebox/config"
	"gorm.io/gorm"
)

type File struct {
	ID        uint           `gorm:"column:id; primarykey" json:"-"`
	Filepath  string         `gorm:"column:filepath; not null" json:"-"` // relative path
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (f *File) GetAbsolutePath() string {
	return path.Join(config.Environment.UploadsPath, f.Filepath)
}

func (f *File) Exists() bool {
	_, err := os.Stat(f.GetAbsolutePath())
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}
