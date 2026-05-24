package notifications

import (
	"sync"

	"gitlab.com/codebox4073715/codebox/db/models"
)

type NotificationMessage struct {
	Type      string
	Event     string
	Workspace *models.Workspace
}

// Notification types
const NotificationTypeWorkspace = "workspace"

// Notification events
const NotificationEventStart = "start"
const NotificationEventStop = "stop"
const NotificationEventRestart = "restart"

type WorkspaceNotificationsHub struct {
	channels []chan NotificationMessage
}

var singletonLock = &sync.Mutex{}
var operationLock = &sync.Mutex{}
var workspaceNotificationsHub *WorkspaceNotificationsHub

/*
Get instance of workspace notifications hub
*/
func GetWorkspaceNotificationsHub() *WorkspaceNotificationsHub {
	if workspaceNotificationsHub == nil {
		singletonLock.Lock()
		defer singletonLock.Unlock()
		if workspaceNotificationsHub == nil {
			workspaceNotificationsHub = &WorkspaceNotificationsHub{
				channels: make([]chan NotificationMessage, 0),
			}
		}
	}

	return workspaceNotificationsHub
}

/*
add notification to queue
*/
func (q *WorkspaceNotificationsHub) SendNotification(
	notification NotificationMessage,
) {
	operationLock.Lock()
	defer operationLock.Unlock()

	for _, ch := range q.channels {
		select {
		case ch <- notification:
		default:
			// client is too slow
		}
	}
}

/*
get channel to listen for notifications, this is used by the websocket handler
to send notifications to the client
*/
func (q *WorkspaceNotificationsHub) GetChannel() chan NotificationMessage {
	operationLock.Lock()
	defer operationLock.Unlock()

	ch := make(chan NotificationMessage, 100)
	q.channels = append(q.channels, ch)
	return ch
}

/*
remove channel from hub, this is used when the websocket connection is closed
to stop sending notifications to the client
*/
func (q *WorkspaceNotificationsHub) RemoveChannel(ch chan NotificationMessage) {
	operationLock.Lock()
	defer operationLock.Unlock()

	for i, c := range q.channels {
		if c == ch {
			q.channels = append(q.channels[:i], q.channels[i+1:]...)
			close(ch)
			break
		}
	}
}

/*
send notification for workspace start
*/
func SendWorkspaceStartNotification(workspace models.Workspace) {
	hub := GetWorkspaceNotificationsHub()
	notification := NotificationMessage{
		Type:      NotificationTypeWorkspace,
		Event:     NotificationEventStart,
		Workspace: &workspace,
	}
	hub.SendNotification(notification)
}

/*
send notification for workspace stop
*/
func SendWorkspaceStopNotification(workspace models.Workspace) {
	hub := GetWorkspaceNotificationsHub()
	notification := NotificationMessage{
		Type:      NotificationTypeWorkspace,
		Event:     NotificationEventStop,
		Workspace: &workspace,
	}
	hub.SendNotification(notification)
}

/*
send notification for workspace restart
*/
func SendWorkspaceRestartNotification(workspace models.Workspace) {
	hub := GetWorkspaceNotificationsHub()
	notification := NotificationMessage{
		Type:      NotificationTypeWorkspace,
		Event:     NotificationEventRestart,
		Workspace: &workspace,
	}
	hub.SendNotification(notification)
}
