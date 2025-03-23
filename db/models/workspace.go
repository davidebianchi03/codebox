package models

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/davidebianchi03/codebox/config"
	"gorm.io/gorm"
)

// workspace status
const WorkspaceStatusStarting = "starting"
const WorkspaceStatusRunning = "running"
const WorkspaceStatusStopping = "stopping"
const WorkspaceStatusStopped = "stopped"
const WorkspaceStatusDeleting = "deleting"
const WorkspaceStatusError = "error"

const (
	WorkspaceConfigSourceGit      = "git"
	WorkspaceConfigSourceTemplate = "template"
)

type Workspace struct {
	ID                   uint                      `gorm:"primarykey" json:"id"`
	Name                 string                    `gorm:"size:255; not null;" json:"name"`
	UserID               uint                      `json:"-"`
	User                 User                      `gorm:"constraint:OnDelete:CASCADE;" json:"user"`
	Status               string                    `gorm:"size:30; not null;" json:"status"`
	Type                 string                    `gorm:"size:255; not null;" json:"type"`
	RunnerID             uint                      `json:"-"`
	Runner               *Runner                   `gorm:"constraint:OnDelete:CASCADE;" json:"runner"`
	ConfigSource         string                    `gorm:"size:20; not null;" json:"config_source"` // template/git
	TemplateVersionID    *uint                     `json:"-"`
	TemplateVersion      *WorkspaceTemplateVersion `gorm:"constraint:OnDelete:CASCADE;" json:"template_version"`
	GitSourceID          *uint                     `json:"-"`
	GitSource            *GitWorkspaceSource       `gorm:"constraint:OnDelete:CASCADE;" json:"git_source"`
	ConfigSourceFilePath string                    `gorm:"type:text;" json:"config_source_file_path"` // name or relative path of the configuration file relative to the template root or repository root folder
	EnvironmentVariables []string                  `gorm:"serializer:json" json:"environment_variables"`
	CreatedAt            time.Time                 `json:"created_at"`
	UpdatedAt            time.Time                 `json:"updated_at"`
	DeletedAt            gorm.DeletedAt            `gorm:"index" json:"-"`
}

func (w *Workspace) GetLogsFilePath() (string, error) {
	workspaceLogsBaseDir := fmt.Sprintf("%s/workspace-logs", config.Environment.UploadsPath)
	err := os.MkdirAll(workspaceLogsBaseDir, 0777)
	if err != nil {
		return "", fmt.Errorf("cannot create path '%s'", workspaceLogsBaseDir)
	}
	return fmt.Sprintf("%s/workspace_%d.log", workspaceLogsBaseDir, w.ID), nil
}

func (w *Workspace) AppendLogs(logs string) error {
	logsFile, err := w.GetLogsFilePath()
	if err != nil {
		return fmt.Errorf("cannot retrieve logs file path, %s", err)
	}
	f, err := os.OpenFile(logsFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("cannot open logs file, %s", err)
	}
	defer f.Close()

	log.SetOutput(f)
	logs = strings.TrimSpace(logs)
	if logs != "" {
		log.Println(logs)
	}
	return nil
}

func (w *Workspace) ClearLogs() error {
	logsFile, err := w.GetLogsFilePath()
	if err != nil {
		return fmt.Errorf("cannot retrieve logs file path, %s", err)
	}
	os.RemoveAll(logsFile)
	return nil
}

func (w *Workspace) RetrieveLogs() (string, error) {
	logsFile, err := w.GetLogsFilePath()
	if err != nil {
		return "", fmt.Errorf("cannot retrieve logs file path, %s", err)
	}

	_, err = os.Stat(logsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	fileContent, err := os.ReadFile(logsFile)
	if err != nil {
		return "", fmt.Errorf("cannot read logs from file, %s", err)
	}
	return string(fileContent), nil
}
