package db

import (
	"fmt"

	"gorm.io/gorm"
)

// workspace type
const (
	WorkspaceContainerTypeDocker = "docker_container"
)

var WorkspaceContainerTypeChoices = [...]string{WorkspaceContainerTypeDocker}

// workspace agent status
const (
	WorkspaceContainerAgentStatusRunning  = "running"
	WorkspaceContainerAgentStatusStarting = "starting"
	WorkspaceContainerAgentStatusError    = "error"
)

var workspaceContainerAgentStatusChoices = [...]string{
	WorkspaceContainerAgentStatusRunning,
	WorkspaceContainerAgentStatusStarting,
	WorkspaceContainerAgentStatusError,
}

// workspace container status
const (
	WorkspaceContainerStatusRunning  = "running"
	WorkspaceContainerStatusStarting = "starting"
	WorkspaceContainerStatusStopped  = "stopped"
	WorkspaceContainerStatusError    = "error"
)

var workspaceContainerStatusChoices = [...]string{
	WorkspaceContainerStatusRunning,
	WorkspaceContainerStatusStarting,
	WorkspaceContainerStatusStopped,
	WorkspaceContainerStatusError,
}

type WorkspaceContainer struct {
	gorm.Model
	Type                       string           `gorm:"size:255; default:docker_container;"`
	Name                       string           `gorm:"size:255;"`
	ContainerUser              string           `gorm:"size:255; default:root;"`
	ContainerStatus            string           `gorm:"size:20; default:starting;"`
	AgentStatus                string           `gorm:"size:20; default:starting;"`
	AgentExternalPort          uint             `gorm:""`
	CanConnectRemoteDeveloping bool             `gorm:"default:false"`
	WorkspacePathInContainer   string           `gorm:"size:1024;"`
	ExternalIPv4               string           `gorm:"size:15;"`
	ForwardedPorts             []*ForwardedPort `gorm:"many2many:workspace_container_forwarded_ports;"`
}

func (wc *WorkspaceContainer) FullClean() (err error) {
	// validate field Type
	if !isItemInArray(wc.Type, WorkspaceContainerTypeChoices[:]) {
		return fmt.Errorf("%s is not a valid value for field 'Type'", wc.Type)
	}

	// validate field ContainerStatus
	if !isItemInArray(wc.ContainerStatus, workspaceContainerStatusChoices[:]) {
		return fmt.Errorf("%s is not a valid value for field 'ContainerStatus'", wc.Type)
	}

	// validate field AgentStatus
	if !isItemInArray(wc.AgentStatus, workspaceContainerAgentStatusChoices[:]) {
		return fmt.Errorf("%s is not a valid value for field 'AgentStatus'", wc.Type)
	}

	return nil
}

func (wc *WorkspaceContainer) BeforeCreate(tx *gorm.DB) (err error) {
	return wc.FullClean()
}

func (wc *WorkspaceContainer) BeforeSave(tx *gorm.DB) (err error) {
	return wc.FullClean()
}
