package models

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/davidebianchi03/codebox/config"
	"gorm.io/gorm"
)

type GitWorkspaceSource struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	RepositoryURL string         `gorm:"type:text;not null;" json:"repository_url"`
	Files         string         `gorm:"type:text;" json:"-"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (gws *GitWorkspaceSource) GetConfigFileAbsPath() (p string, err error) {
	if gws.ID <= 0 {
		return "", errors.New("object does not exist")
	}
	folder := path.Join(config.Environment.UploadsPath, "git-sources")
	if err = os.MkdirAll(folder, 0777); err != nil {
		return "", err
	}
	p = path.Join(folder, fmt.Sprintf("%d.tar.gz", gws.ID))
	return p, nil
}
