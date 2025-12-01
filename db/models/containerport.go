package models

import (
	"time"

	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gorm.io/gorm"
)

type WorkspaceContainerPort struct {
	ID          uint               `gorm:"primarykey" json:"-"`
	ContainerID uint               `gorm:"column:container_id;" json:"-"`
	Container   WorkspaceContainer `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ServiceName string             `gorm:"column:service_name; size:255; not null;" json:"service_name"`
	PortNumber  uint               `gorm:"column:port_number; not null;" json:"port_number"`
	Public      bool               `gorm:"column:public; default:false;" json:"public"`
	CreatedAt   time.Time          `gorm:"column:created_at;" json:"created_at"`
	UpdatedAt   time.Time          `gorm:"column:updated_at;" json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `gorm:"index" json:"-"`
}

func ListContainerPortsByWorkspaceContainer(container WorkspaceContainer) ([]WorkspaceContainerPort, error) {
	var ports []WorkspaceContainerPort
	result := dbconn.DB.
		Preload("Container").
		Preload("Container.Workspace").
		Where("container_id = ?", container.ID).
		Find(&ports)

	if result.Error != nil {
		return nil, result.Error
	}
	return ports, nil
}

func RetrieveContainerPortByServiceName(
	container WorkspaceContainer,
	serviceName string,
) (*WorkspaceContainerPort, error) {
	var port WorkspaceContainerPort
	result := dbconn.DB.
		Preload("Container").
		Preload("Container.Workspace").
		Where("container_id = ? AND service_name = ?", container.ID, serviceName).
		First(&port)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // No port found
		}
		return nil, result.Error
	}
	return &port, nil
}

func RetrieveContainerPortByPortNumber(
	container WorkspaceContainer,
	portNumber uint,
) (*WorkspaceContainerPort, error) {
	var port WorkspaceContainerPort
	result := dbconn.DB.
		Preload("Container").
		Preload("Container.Workspace").
		Where("container_id = ? AND port_number = ?", container.ID, portNumber).
		First(&port)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // No port found
		}
		return nil, result.Error
	}
	return &port, nil
}

func CreateContainerPort(
	container WorkspaceContainer,
	serviceName string,
	portNumber uint,
	public bool,
) (*WorkspaceContainerPort, error) {
	port := &WorkspaceContainerPort{
		ContainerID: container.ID,
		ServiceName: serviceName,
		PortNumber:  portNumber,
		Public:      public,
	}

	result := dbconn.DB.Create(port)
	if result.Error != nil {
		return nil, result.Error
	}
	return port, nil
}

func UpdateContainerPort(
	port *WorkspaceContainerPort,
	serviceName string,
	portNumber uint,
	public bool,
) (*WorkspaceContainerPort, error) {
	port.ServiceName = serviceName
	port.PortNumber = portNumber
	port.Public = public

	result := dbconn.DB.Save(port)
	if result.Error != nil {
		return nil, result.Error
	}
	return port, nil
}

func DeleteContainerPort(port *WorkspaceContainerPort) error {
	result := dbconn.DB.Unscoped().Delete(port)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
