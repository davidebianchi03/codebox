package models

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gitlab.com/codebox4073715/codebox/config"
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
	Name                 string                    `gorm:"column:name; size:255; not null;" json:"name"`
	UserID               uint                      `gorm:"column:user_id;" json:"-"`
	User                 *User                     `gorm:"constraint:OnDelete:CASCADE;" json:"user"`
	Status               string                    `gorm:"column:status; size:30; not null;" json:"status"`
	Type                 string                    `gorm:"column:type; size:255; not null;" json:"type"`
	RunnerID             uint                      `gorm:"column:runner_id;" json:"-"`
	Runner               *Runner                   `gorm:"constraint:OnDelete:CASCADE;" json:"runner"`
	ConfigSource         string                    `gorm:"column:config_source; size:20; not null;" json:"config_source"` // template/git
	TemplateVersionID    *uint                     `gorm:"column:template_version_id;" json:"-"`
	TemplateVersion      *WorkspaceTemplateVersion `gorm:"constraint:OnDelete:CASCADE;" json:"template_version"`
	GitSourceID          *uint                     `gorm:"column:git_source_id;" json:"-"`
	GitSource            *GitWorkspaceSource       `gorm:"constraint:OnDelete:CASCADE;" json:"git_source"`
	EnvironmentVariables []string                  `gorm:"column:environment_variables; serializer:json" json:"environment_variables"`
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

/*
Retrieve list of environement variables to exported by default to workspace
This variables are informations related to workspace such as workspace name,
owner email, first name and last name
*/
func (w *Workspace) GetDefaultEnvironmentVariables() []string {
	return []string{
		fmt.Sprintf("CODEBOX_WORKSPACE_ID=%d", w.ID),
		fmt.Sprintf("CODEBOX_WORKSPACE_NAME=%s", w.Name),
		fmt.Sprintf("CODEBOX_WORKSPACE_OWNER_EMAIL=%s", w.User.Email),
		fmt.Sprintf("CODEBOX_WORKSPACE_OWNER_FIRST_NAME=%s", w.User.FirstName),
		fmt.Sprintf("CODEBOX_WORKSPACE_OWNER_LAST_NAME=%s", w.User.LastName),
		fmt.Sprintf("CODEBOX_WORKSPACE_RUNNER_ID=%d", w.Runner.ID),
		fmt.Sprintf("CODEBOX_WORKSPACE_RUNNER_NAME=%s", w.Runner.Name),
	}
}
