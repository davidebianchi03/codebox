package serializers

import "gitlab.com/codebox4073715/codebox/httpserver/notifications"

type NotificationMessageSerializers struct {
	Type      string               `json:"type"`
	Event     string               `json:"event"`
	Workspace *WorkspaceSerializer `json:"workspace,omitempty"`
}

func LoadNotificationMessageSerializer(
	notificationMessage notifications.NotificationMessage,
) NotificationMessageSerializers {
	return NotificationMessageSerializers{
		Type:      notificationMessage.Type,
		Event:     notificationMessage.Event,
		Workspace: LoadWorkspaceSerializer(notificationMessage.Workspace),
	}
}
