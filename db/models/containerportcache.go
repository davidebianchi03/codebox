package models

import (
	"time"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

/*
ContainerPortBackup is a backup of the ContainerPort model, used to
save exposed ports and automatically restore them when a container is recreated.
*/

type ContainerPortBackup struct {
	ID            uint           `gorm:"primarykey" json:"-"`
	WorkspaceID   uint           `gorm:"column:workspace_id;" json:"-"`
	Workspace     Workspace      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ContainerName string         `gorm:"column:container_name; size:255; not null" json:"container_name"`
	ServiceName   string         `gorm:"column:service_name; size:255; not null;" json:"service_name"`
	PortNumber    uint           `gorm:"column:port_number; not null;" json:"port_number"`
	Public        bool           `gorm:"column:public; default:false;" json:"public"`
	CreatedAt     time.Time      `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

/*
ListContainerPortBackupByWorkspace retrieves all port backup entries for a given workspace.
*/
func ListContainerPortBackupByWorkspace(workspace Workspace) ([]ContainerPortBackup, error) {
	var containers []ContainerPortBackup
	result := dbconn.DB.
		Preload("Workspace").
		Where("workspace_id = ?", workspace.ID).
		Find(&containers)

	if result.Error != nil {
		return nil, result.Error
	}
	return containers, nil
}

/*
Add a new ContainerPortBackup entry to the database.
*/
func CreateContainerPortBackup(
	container WorkspaceContainer,
	serviceName string,
	portNumber uint,
	public bool,
) (*ContainerPortBackup, error) {
	port := &ContainerPortBackup{
		WorkspaceID:   container.WorkspaceID,
		ContainerName: container.ContainerName,
		ServiceName:   serviceName,
		PortNumber:    portNumber,
		Public:        public,
	}

	result := dbconn.DB.Create(port)
	if result.Error != nil {
		return nil, result.Error
	}
	return port, nil
}

/*
DeleteWorkspaceContainerPortsBackups deletes all port backup
entries for a given workspace.
*/
func DeleteWorkspaceContainerPortsBackups(workspace Workspace) error {
	containers, err := ListContainerPortBackupByWorkspace(workspace)
	if err != nil {
		return err
	}

	for _, container := range containers {
		if err := dbconn.DB.Unscoped().Delete(&container).Error; err != nil {
			return err
		}
	}
	return nil
}
