package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	TypeDevcontainer = "devcontainer"
)

var WorkspaceTypeChoices = [...]string{
	TypeDevcontainer,
}

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
	Name                        string                 `gorm:"column:name; size:100; not null;"`
	Status                      string                 `gorm:"column:status; size:40; not null;default:creating;"`
	OwnerId                     uint                   `gorm:"column:owner_id;"`
	Owner                       User                   `gorm:"foreignKey:OwnerId; references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null;"`
	Type                        string                 `gorm:"column:type; size:100; not null; default:devcontainer;"`
	GitRepoUrl                  string                 `gorm:"column:git_repo_url; size:1024;"`
	GitRepoConfigurationFolder  string                 `gorm:"column:git_repo_configuration_folder; size:255;"`
	LastActivityOn              time.Time              `gorm:"column:last_activity_on;"`
	LastStartOn                 time.Time              `gorm:"column:last_start_on;"`
	Logs                        string                 `gorm:"column:logs;"`
	WorkspaceConfigurationFiles string                 `gorm:"column:workspace_configuration_files; size:1024;"`
	CustomConfig                map[string]interface{} `gorm:"column:custom_config; serializer:json"` // configurazione che varia dal tipo di workspace
}

func (w *Workspace) FullClean() (err error) {
	if w.Type != TypeDevcontainer {
		return fmt.Errorf("workspace type: '%s' is not supported", w.Type)
	}

	if !IsItemInArray(w.Status, workspaceStatusChoices[:]) {
		return fmt.Errorf("invalid workspace status: '%s'", w.Status)
	}
	return nil
}

func (w *Workspace) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func (w *Workspace) BeforeSave(tx *gorm.DB) (err error) {
	return w.FullClean()
}
