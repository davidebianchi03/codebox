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
	Type                       string           `gorm:"column:type; size:255; default:docker_container;"`
	Name                       string           `gorm:"column:name; size:255;"`
	ContainerUser              string           `gorm:"column:container_user; size:255; default:root;"`
	ContainerStatus            string           `gorm:"column:container_status; size:20; default:starting;"`
	AgentStatus                string           `gorm:"column:agent_status; size:20; default:starting;"`
	AgentExternalPort          uint             `gorm:"column:agent_external_port;"`
	CanConnectRemoteDeveloping bool             `gorm:"column:can_connect_remote_developing; default:false"`
	WorkspacePathInContainer   string           `gorm:"column:workspace_path_in_container; size:1024;"`
	ExternalIPv4               string           `gorm:"column:external_ipv4; size:15;"`
	ForwardedPorts             []*ForwardedPort `gorm:"many2many:workspace_container_forwarded_ports;"`
}

func (wc *WorkspaceContainer) FullClean() (err error) {
	// validate field Type
	if !IsItemInArray(wc.Type, WorkspaceContainerTypeChoices[:]) {
		return fmt.Errorf("%s is not a valid value for field 'Type'", wc.Type)
	}

	// validate field ContainerStatus
	if !IsItemInArray(wc.ContainerStatus, workspaceContainerStatusChoices[:]) {
		return fmt.Errorf("%s is not a valid value for field 'ContainerStatus'", wc.Type)
	}

	// validate field AgentStatus
	if !IsItemInArray(wc.AgentStatus, workspaceContainerAgentStatusChoices[:]) {
		return fmt.Errorf("%s is not a valid value for field 'AgentStatus'", wc.Type)
	}

	return nil
}

func (wc *WorkspaceContainer) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

func (wc *WorkspaceContainer) BeforeSave(tx *gorm.DB) (err error) {
	return wc.FullClean()
}
