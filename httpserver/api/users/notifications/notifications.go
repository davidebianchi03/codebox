package notifications

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.com/codebox4073715/codebox/httpserver/api/users/serializers"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
	"gitlab.com/codebox4073715/codebox/httpserver/notifications"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Accepting all requests
	},
}

// HandleWorkspaceNotifications godoc
// @Summary WebSocket connection for workspace notifications
// @Schemes ws wss
// @Description Establish WebSocket connection to receive real-time notifications about workspace status changes and template version releases
// @Tags Notifications
// @Accept json
// @Produce json
// @Success 101 {object} string "WebSocket upgrade"
// @Failure 400 {object} string "Cannot upgrade connection"
// @Failure 500 {object} string "Internal server error"
// @Router /api/v1/notifications [get]
func HandleWorkspaceNotifications(c *gin.Context) {
	currentUser, err := utils.GetUserFromContext(c)
	if err != nil {
		// TODO: log error
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	if currentUser.ID == 0 {
		// TODO: log error
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusBadRequest,
			"cannot upgrade ws connection",
		)
		return
	}
	defer wsConn.Close()

	hub := notifications.GetWorkspaceNotificationsHub()
	channel := hub.GetChannel()
	defer hub.RemoveChannel(channel)

	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	for {
		// send ping message to keep the connection alive
		select {
		case <-pingTicker.C:
			if err := wsConn.WriteMessage(
				websocket.PingMessage,
				[]byte("ping")); err != nil {
				return
			}
		default:
			// no ping to send
		}

		select {
		case notification := <-channel:
			// only send notifications related to the user's workspaces
			if notification.Workspace.UserID == currentUser.ID {
				if err := wsConn.WriteJSON(
					serializers.LoadNotificationMessageSerializer(notification),
				); err != nil {
					return
				}
			}
		default:
			// no notifications to send
		}
	}
}
