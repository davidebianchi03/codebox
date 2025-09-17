package models

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gitlab.com/codebox4073715/codebox/config"
	"gorm.io/gorm"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
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
		fmt.Sprintf("CODEBOX_WORKSPACE_NAME=%s", strings.ToLower(w.Name)),
		fmt.Sprintf("CODEBOX_WORKSPACE_OWNER_EMAIL=%s", w.User.Email),
		fmt.Sprintf("CODEBOX_WORKSPACE_OWNER_FIRST_NAME=%s", strings.ToLower(w.User.FirstName)),
		fmt.Sprintf("CODEBOX_WORKSPACE_OWNER_LAST_NAME=%s", strings.ToLower(w.User.LastName)),
		fmt.Sprintf("CODEBOX_WORKSPACE_RUNNER_ID=%d", w.Runner.ID),
		fmt.Sprintf("CODEBOX_WORKSPACE_RUNNER_NAME=%s", strings.ToLower(w.Runner.Name)),
	}
}

/*
Filter workspaces by owner
*/
func ListUserWorkspaces(user User) ([]Workspace, error) {
	workspaces := []Workspace{}
	r := dbconn.DB.
		Preload("GitSource").
		Preload("TemplateVersion").
		Preload("Runner").
		Preload("User").
		Find(
			&workspaces,
			map[string]interface{}{
				"user_id": user.ID,
			},
		)

	if r.Error != nil {
		return []Workspace{}, r.Error
	}

	return workspaces, nil
}

/*
Retrieve a workspace by workspace id and owner.
If workspace does not exist return nil.
*/
func RetrieveWorkspaceByUserAndId(user User, id uint) (*Workspace, error) {
	workspace := Workspace{}
	r := dbconn.DB.
		Preload("GitSource").
		Preload("TemplateVersion").
		Preload("Runner").
		Preload("User").
		Find(
			&workspace,
			map[string]interface{}{
				"ID":      id,
				"user_id": user.ID,
			},
		)

	if r.Error != nil {
		if r.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, r.Error
	}

	if r.RowsAffected == 1 {
		return &workspace, nil
	}

	return nil, nil
}

/*
Create a new workspace, this function will add new row to the database.
It will not start the workspace, it will only create the database entry.
*/
func CreateWorkspace(
	name string,
	user *User,
	workspaceType string,
	runner *Runner,
	configSource string,
	templateVersion *WorkspaceTemplateVersion,
	gitSource *GitWorkspaceSource,
	environmentVariables []string,
) (*Workspace, error) {

	var templateVersionID *uint
	if templateVersion != nil {
		templateVersionID = &templateVersion.ID
	} else {
		templateVersionID = nil
	}

	var gitSourceID *uint
	if gitSource != nil {
		gitSourceID = &gitSource.ID
	} else {
		gitSourceID = nil
	}

	workspace := Workspace{
		Name:                 name,
		UserID:               user.ID,
		User:                 user,
		Status:               WorkspaceStatusStarting,
		Type:                 workspaceType,
		RunnerID:             runner.ID,
		Runner:               runner,
		ConfigSource:         configSource,
		TemplateVersionID:    templateVersionID,
		TemplateVersion:      templateVersion,
		GitSourceID:          gitSourceID,
		GitSource:            gitSource,
		EnvironmentVariables: environmentVariables,
	}

	r := dbconn.DB.Create(&workspace)
	if r.Error != nil {
		return nil, r.Error
	}
	return &workspace, nil
}

/*
Update a workspace
*/
func UpdateWorkspace(
	workspace *Workspace,
	name string,
	status string,
	runner *Runner,
	configSource string,
	templateVersion *WorkspaceTemplateVersion,
	gitSource *GitWorkspaceSource,
	environmentVariables []string,
) (*Workspace, error) {

	workspace.Name = name
	workspace.Status = status
	workspace.RunnerID = runner.ID
	workspace.Runner = runner
	workspace.ConfigSource = configSource
	if templateVersion != nil {
		workspace.TemplateVersionID = &templateVersion.ID
	} else {
		workspace.TemplateVersionID = nil
	}
	workspace.TemplateVersion = templateVersion
	if gitSource != nil {
		workspace.GitSourceID = &gitSource.ID
	} else {
		workspace.GitSourceID = nil
	}
	workspace.GitSource = gitSource
	workspace.EnvironmentVariables = environmentVariables

	if err := dbconn.DB.Save(&workspace).Error; err != nil {
		return nil, err
	}

	return workspace, nil
}

/*
CountAllOnlineWorkspaces counts the number of workspaces that are currently running.
*/
func CountAllOnlineWorkspaces() (int64, error) {
	var count int64
	if err := dbconn.DB.Model(&Workspace{}).
		Where("status = ?", WorkspaceStatusRunning).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
