package models

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/davidebianchi03/codebox/config"
	"gorm.io/gorm"
)

// workspace status
const (
	WorkspaceStatusCreating = "creating"
	WorkspaceStatusRunning  = "running"
	WorkspaceStatusStopping = "stopping"
	WorkspaceStatusStopped  = "stopped"
	WorkspaceStatusStarting = "starting"
	WorkspaceStatusDeleting = "deleting"
	WorkspaceStatusError    = "error"
)

type Workspace struct {
	gorm.Model
	ID                   uint   `gorm:"primarykey"`
	Name                 string `gorm:"size:255; not null;"`
	UserID               uint
	User                 User   `gorm:"constraint:OnDelete:CASCADE;"`
	Status               string `gorm:"size:30; not null;"`
	Type                 string `gorm:"size:255; not null;"`
	RunnerID             uint
	Runner               Runner `gorm:"constraint:OnDelete:CASCADE;"`
	ConfigSource         string `gorm:"size:20; not null;"` // template/git
	TemplateVersionID    uint
	TemplateVersion      WorkspaceTemplateVersion `gorm:"constraint:OnDelete:CASCADE;"`
	GitSourceID          uint
	GitSource            GitWorkspaceSource `gorm:"constraint:OnDelete:CASCADE;"`
	EnvironmentVariables []string           `gorm:"serializer:json"`
}

func (w *Workspace) GetLogsFilePath() (string, error) {
	workspaceLogsBaseDir := fmt.Sprintf("%s/workspace-logs", config.Environment.UploadsPath)

	info, err := os.Stat(workspaceLogsBaseDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(workspaceLogsBaseDir, 0777)
			if err != nil {
				return "", fmt.Errorf("cannot create path '%s'", workspaceLogsBaseDir)
			}
		}
		return "", fmt.Errorf("unknown error")
	}

	if !info.IsDir() {
		return "", fmt.Errorf("logs path is not a directory '%s'", workspaceLogsBaseDir)
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
	fileContent, err := os.ReadFile(logsFile)
	if err != nil {
		return "", fmt.Errorf("cannot read logs from file, %s", err)
	}
	return string(fileContent), nil
}
