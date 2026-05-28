package notifications_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver"
	"gitlab.com/codebox4073715/codebox/httpserver/notifications"
	"gitlab.com/codebox4073715/codebox/testutils"
)

/*
Test notifications endpoint
*/
func TestWorkspaceNotifications(t *testing.T) {
	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
		router := httpserver.SetupRouter()
		server := httptest.NewServer(router)
		defer server.Close()

		user, err := models.RetrieveUserByEmail("user1@user.com")
		assert.NoError(t, err)
		assert.NotNil(t, user)

		runners, err := models.ListRunners(1, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, runners)

		gitSource, err := models.CreateGitWorkspaceSource(
			"https://git.example.com/test",
			"refs/head/master",
			"docker-compose.yml",
		)
		assert.NoError(t, err)

		workspace, err := models.CreateWorkspace(
			"test",
			user,
			"docker_compose",
			&runners[0],
			models.WorkspaceConfigSourceGit,
			nil,
			gitSource,
			[]string{"VAR1=value1"},
		)
		assert.NoError(t, err)

		token, err := models.CreateToken(*user, time.Hour*24*20)
		assert.NoError(t, err)

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/notifications"

		header := http.Header{}
		header.Set("Authorization", "Bearer "+token.Token)

		wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
		assert.NoError(t, err)
		defer wsConn.Close()

		messages := make(chan string, 10)
		done := make(chan struct{})

		// reader goroutine
		go func() {
			defer close(done)

			for {
				wsConn.SetReadDeadline(time.Now().Add(10 * time.Second))

				_, data, err := wsConn.ReadMessage()
				if err != nil {
					return
				}

				messages <- string(data)
			}
		}()

		// trigger event AFTER websocket is ready
		time.Sleep(200 * time.Millisecond) // fallback minimo (vedi nota sotto)
		notifications.SendWorkspaceStartNotification(*workspace)

		// wait for message deterministically
		timeout := time.After(10 * time.Second)

		for {
			select {
			case msg := <-messages:
				t.Logf("received: %s", msg)

				var notification map[string]interface{}
				err := json.Unmarshal([]byte(msg), &notification)
				assert.NoError(t, err)

				if notification["type"] == "workspace" && notification["event"] == "start" {
					assert.NotNil(t, notification["workspace"])
					return
				}

			case <-done:
				t.Fatal("websocket closed unexpectedly")

			case <-timeout:
				t.Fatal("timeout waiting for workspace notification")
			}
		}
	})
}

// /*
// Test notifications endpoint, try to start a workspace
// then verify that another user does not receive the notification
// */
// func TestWorkspaceNotificationsDifferentUser(t *testing.T) {
// 	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
// 		router := httpserver.SetupRouter()
// 		server := httptest.NewServer(router)
// 		defer server.Close()

// 		user1, err := models.RetrieveUserByEmail("user1@user.com")
// 		if err != nil || user1 == nil {
// 			t.Errorf("Failed to retrieve user: '%s'\n", err)
// 			t.FailNow()
// 		}

// 		user2, err := models.RetrieveUserByEmail("user2@user.com")
// 		if err != nil || user2 == nil {
// 			t.Errorf("Failed to retrieve user: '%s'\n", err)
// 			t.FailNow()
// 		}

// 		runners, err := models.ListRunners(1, 0)
// 		if err != nil || len(runners) == 0 {
// 			t.Errorf("Failed to retrieve runner: '%s'\n", err)
// 			t.FailNow()
// 		}

// 		// create a workspace
// 		gitSource, err := models.CreateGitWorkspaceSource(
// 			"https://git.example.com/test",
// 			"refs/head/master",
// 			"docker-compose.yml",
// 		)
// 		assert.Nil(t, err)

// 		workspace, err := models.CreateWorkspace(
// 			"test",
// 			user1,
// 			"docker_compose",
// 			&runners[0],
// 			models.WorkspaceConfigSourceGit,
// 			nil,
// 			gitSource,
// 			[]string{"VAR1=value1"},
// 		)
// 		assert.Nil(t, err)

// 		// list notifications with user2 token
// 		token, err := models.CreateToken(*user2, time.Duration(time.Hour*24*20))
// 		if err != nil {
// 			t.Errorf("Failed to create token: '%s'\n", err)
// 			t.FailNow()
// 		}

// 		// connect to WebSocket
// 		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/notifications"
// 		header := http.Header{}
// 		header.Set("Authorization", "Bearer "+token.Token)
// 		wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
// 		if err != nil {
// 			t.Errorf("Failed to connect to WebSocket: '%s'\n", err)
// 			t.FailNow()
// 		}
// 		defer wsConn.Close()

// 		messages := make(chan string, 10)
// 		done := make(chan bool, 1)

// 		// read messages in a separate goroutine
// 		go func() {
// 			for {
// 				wsConn.SetReadDeadline(time.Now().Add(5 * time.Second))
// 				_, data, err := wsConn.ReadMessage()
// 				if err != nil {
// 					done <- true
// 					return
// 				}
// 				messages <- string(data)
// 			}
// 		}()

// 		// sleep for a while
// 		time.Sleep(100 * time.Millisecond)

// 		timeout := time.NewTimer(2 * time.Second)
// 		var notification map[string]interface{}

// 		// send workspace start notification
// 		notifications.SendWorkspaceStartNotification(*workspace)

// 		for {
// 			select {
// 			case msg := <-messages:
// 				t.Logf("Received notification: %s\n", msg)
// 				if err := json.Unmarshal([]byte(msg), &notification); err != nil {
// 					t.Logf("Failed to parse notification: %s\n", err)
// 				}
// 			case <-done:
// 				t.Logf("WebSocket connection closed")
// 				return
// 			case <-timeout.C:
// 				break
// 			}
// 			if notification != nil {
// 				break
// 			}
// 		}

// 		timeout.Stop()
// 		assert.Nil(t, notification)
// 	})
// }

// /*
// Test notifications endpoint without authentication
// */
// func TestWorkspaceNotificationsUnauthorized(t *testing.T) {
// 	testutils.WithSetupAndTearDownTestEnvironment(t, func(t *testing.T) {
// 		router := httpserver.SetupRouter()
// 		server := httptest.NewServer(router)
// 		defer server.Close()

// 		// connect to WebSocket
// 		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/notifications"
// 		header := http.Header{}
// 		_, resp, err := websocket.DefaultDialer.Dial(wsURL, header)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
// 	})
// }
