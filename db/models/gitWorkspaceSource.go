package models

import (
	"errors"
	"path"
	"strconv"
	"time"

	"github.com/davidebianchi03/codebox/config"
	"gorm.io/gorm"
)

type GitWorkspaceSource struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	RepositoryURL string         `gorm:"size:1024;not null;" json:"repository_url"`
	Files         string         `gorm:"size:1024;" json:"-"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (gws *GitWorkspaceSource) GetConfigFileAbsPath() (p string, err error) {
	if gws.ID <= 0 || gws.Files == "" {
		return "", errors.New("object does not exist")
	}
	p = path.Join(config.Environment.UploadsPath, "git-sources", strconv.Itoa(int(gws.ID)))
	return p, nil
}
