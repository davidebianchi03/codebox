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
const NotificationEventRunning = "running"
const NotificationEventStopped = "stopped"

type ClientChannel struct {
	UserID int
	Ch     chan NotificationMessage
}

type WorkspaceNotificationsHub struct {
	mu       sync.RWMutex
	channels map[chan NotificationMessage]*ClientChannel
}

var singletonLock = &sync.Mutex{}
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
				channels: make(map[chan NotificationMessage]*ClientChannel),
			}
		}
	}

	return workspaceNotificationsHub
}

/*
add notification to queue - optimized to only send to affected users
*/
func (q *WorkspaceNotificationsHub) SendNotification(
	notification NotificationMessage,
) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	// Only send to clients of the workspace owner
	targetUserID := notification.Workspace.UserID

	for _, client := range q.channels {
		if uint(client.UserID) == targetUserID {
			select {
			case client.Ch <- notification:
			default:
				// client is too slow, drop notification to prevent blocking
			}
		}
	}
}

/*
get channel to listen for notifications, this is used by the websocket handler
to send notifications to the client
*/
func (q *WorkspaceNotificationsHub) GetChannel(userID int) chan NotificationMessage {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Use buffered channel with smaller buffer to prevent slowloris
	ch := make(chan NotificationMessage, 32)
	q.channels[ch] = &ClientChannel{
		UserID: userID,
		Ch:     ch,
	}
	return ch
}

/*
remove channel from hub, this is used when the websocket connection is closed
to stop sending notifications to the client
*/
func (q *WorkspaceNotificationsHub) RemoveChannel(ch chan NotificationMessage) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, exists := q.channels[ch]; exists {
		delete(q.channels, ch)
		close(ch)
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

/*
send notification for workspace running
*/
func SendWorkspaceRunningNotification(workspace models.Workspace) {
	hub := GetWorkspaceNotificationsHub()
	notification := NotificationMessage{
		Type:      NotificationTypeWorkspace,
		Event:     NotificationEventRunning,
		Workspace: &workspace,
	}
	hub.SendNotification(notification)
}

/*
send notification for workspace stopped
*/
func SendWorkspaceStoppedNotification(workspace models.Workspace) {
	hub := GetWorkspaceNotificationsHub()
	notification := NotificationMessage{
		Type:      NotificationTypeWorkspace,
		Event:     NotificationEventStopped,
		Workspace: &workspace,
	}
	hub.SendNotification(notification)
}
