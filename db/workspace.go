package db

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"codebox.com/env"
	"gorm.io/gorm"
)

const (
	WorkspaceTypeDevcontainer = "devcontainer"
)

var WorkspaceTypeChoices = [...]string{
	WorkspaceTypeDevcontainer,
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
	Name                       string    `gorm:"column:name; size:100; not null;"`
	Status                     string    `gorm:"column:status; size:40; not null;default:creating;"`
	OwnerId                    uint      `gorm:"column:owner_id;"`
	Owner                      User      `gorm:"foreignKey:OwnerId; references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null;"`
	Type                       string    `gorm:"column:type; size:100; not null; default:devcontainer;"`
	GitRepoUrl                 string    `gorm:"column:git_repo_url; size:1024;"`
	GitRepoConfigurationFolder string    `gorm:"column:git_repo_configuration_folder; size:255;"`
	LastActivityOn             time.Time `gorm:"column:last_activity_on;"`
	LastStartOn                time.Time `gorm:"column:last_start_on;"`
	// Logs                        string                 `gorm:"column:logs;"`
	WorkspaceConfigurationFiles string                 `gorm:"column:workspace_configuration_files; size:1024;"`
	CustomConfig                map[string]interface{} `gorm:"column:custom_config; serializer:json"` // configurazione che varia dal tipo di workspace
}

func (w *Workspace) FullClean() (err error) {
	if w.Type != WorkspaceTypeDevcontainer {
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

func (w *Workspace) GetConfigFilePath() (string, error) {
	path := fmt.Sprintf("%s/workspace-configs", env.CodeBoxEnv.UploadsPath)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(path, os.FileMode(0777))
		} else {
			return "", fmt.Errorf("unknown error: %s", err)
		}
	}
	os.Chown(path, -1, -1)
	return fmt.Sprintf("%s/workspace_%d.tar.gz", path, w.ID), nil
}

func (w *Workspace) GetLogsFilePath() (string, error) {
	workspaceLogsBaseDir := fmt.Sprintf("%s/workspace-logs", env.CodeBoxEnv.UploadsPath)

	info, err := os.Stat(workspaceLogsBaseDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(workspaceLogsBaseDir, 0777)
			if err != nil {
				return "", fmt.Errorf("Cannot create path '%s'", workspaceLogsBaseDir)
			}
		}
		return "", fmt.Errorf("Unknown error")
	}

	if !info.IsDir() {
		return "", fmt.Errorf("Logs path is not a directory '%s'", workspaceLogsBaseDir)
	}

	return fmt.Sprintf("%s/workspace_%d.log", workspaceLogsBaseDir, w.ID), nil
}

func (w *Workspace) AppendLogs(logs string) error {
	logsFile, err := w.GetLogsFilePath()
	if err != nil {
		return fmt.Errorf("Cannot retrieve logs file path, %s", err)
	}
	f, err := os.OpenFile(logsFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("Cannot open logs file, %s", err)
	}
	defer f.Close()

	log.SetOutput(f)
	logs = strings.TrimSpace(logs)
	if logs != "" {
		log.Print(logs)
	}
	return nil
}

func (w *Workspace) ClearLogs() error {
	logsFile, err := w.GetLogsFilePath()
	if err != nil {
		return fmt.Errorf("Cannot retrieve logs file path, %s", err)
	}
	err = os.WriteFile(logsFile, []byte(""), os.FileMode(777))
	if err != nil {
		return fmt.Errorf("Cannot clear logs file, %s", err)
	}
	return nil
}

func (w *Workspace) RetrieveLogs() (string, error) {
	logsFile, err := w.GetLogsFilePath()
	if err != nil {
		return "", fmt.Errorf("Cannot retrieve logs file path, %s", err)
	}
	fileContent, err := os.ReadFile(logsFile)
	if err != nil {
		return "", fmt.Errorf("Cannot read logs from file, %s", err)
	}
	return string(fileContent), nil
}
