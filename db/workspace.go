package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	TypeDevcontainer = "devcontainer"
)

// workspace status
const (
	WorkspaceStatusCreating = "creating"
	WorkspaceStatusRunning  = "running"
	WorkspaceStatusStopping = "stopping"
	WorkspaceStatusStopped  = "stopped"
	WorkspaceStatusStarting = "starting"
	WorkspaceStatusError    = "error"
)

var workspaceStatusChoices = [...]string{
	WorkspaceStatusCreating,
	WorkspaceStatusRunning,
	WorkspaceStatusStopping,
	WorkspaceStatusStopped,
	WorkspaceStatusStarting,
	WorkspaceStatusError,
}

type Workspace struct {
	gorm.Model
	Name                        string                 `gorm:"size:100; not null;"`
	Status                      string                 `gorm:"size:40; not null;default:creating;"`
	OwnerId                     uint                   `gorm:""`
	Owner                       User                   `gorm:"foreignKey:OwnerId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null;"`
	Type                        string                 `gorm:"size:100; not null; default:devcontainer;"`
	GitRepoUrl                  string                 `gorm:"size:1024;"`
	GitRepoConfigurationFolder  string                 `gorm:"size:255;"`
	CreatedOn                   time.Time              `gorm:""`
	LastActivityOn              time.Time              `gorm:""`
	LastStartOn                 time.Time              `gorm:""`
	LogsFile                    string                 `gorm:"size:1024;"`
	WorkspaceConfigurationFiles string                 `gorm:"size:1024;"`
	CustomConfig                map[string]interface{} `gorm:"serializer:json"` // configurazione che varia dal tipo di workspace
}

func (w *Workspace) FullClean() (err error) {
	if w.Type != TypeDevcontainer {
		return fmt.Errorf("workspace type: '%s' is not supported", w.Type)
	}

	if isItemInArray(w.Status, workspaceStatusChoices[:]) {
		return fmt.Errorf("invalid workspace status: '%s'", w.Status)
	}
	return nil
}

func (w *Workspace) BeforeCreate(tx *gorm.DB) (err error) {
	w.CreatedOn = time.Now()
	return w.FullClean()
}

func (w *Workspace) BeforeSave(tx *gorm.DB) (err error) {
	return w.FullClean()
}
