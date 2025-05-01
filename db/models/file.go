package models

import (
	"os"
	"path"
	"path/filepath"
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

// get the absolute location of the file
// create parent directory if it does not exists
func (f *File) GetAbsolutePath() string {
	fp := path.Join(config.Environment.UploadsPath, f.Filepath)

	// create parent folder if it does not exists
	parentFolder := filepath.Dir(fp)
	os.MkdirAll(parentFolder, 0700)

	return fp
}

func (f *File) Exists() bool {
	_, err := os.Stat(f.GetAbsolutePath())
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}
