package runners

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/httpserver/api/utils"
	"golang.org/x/crypto/ssh"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Accepting all requests
	},
}

const GitSSHPingInterval = 1 * time.Second
const GitSSHCloseDelay = 100 * time.Millisecond

// RunnerGitSSH godoc
// @Summary Handle ws connection to perform git pulls/pushs over ssh
// @Schemes
// @Description Handle ws connection to perform git pulls/pushs over ssh
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200
// @Error 400
func HandleRunnerGitSSH(c *gin.Context) {
	runnerId, _ := utils.GetUIntParamFromContext(c, "runnerId")
	workspaceId, _ := utils.GetUIntParamFromContext(c, "workspaceId")

	// Container name is unused for now but we might need it in the future
	// to do more checks or logging. For now we can't use it because the container
	// instance may not exists when the ws connection is established.
	// Containers are authenticated by runners
	containerName, _ := c.Params.Get("containerName")
	_ = containerName

	runner, err := models.RetrieveRunnerByID(runnerId)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	// parse url
	q := c.Request.URL.Query()
	sshHost := q.Get("ssh_host")
	sshUser := q.Get("ssh_user")
	sshCmd := q.Get("ssh_cmd")
	repoPath := q.Get("repo_path")
	if sshHost == "" || sshCmd == "" || repoPath == "" || sshUser == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing or invalid parameter")
		return
	}

	// TODO: check that sshcmd is a valid git ssh command

	// retrieve workspace
	workspace, err := models.RetrieveWorkspaceById(workspaceId)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	if workspace == nil {
		utils.ErrorResponse(
			c,
			http.StatusNotFound,
			"workspace not found",
		)
		return
	}

	// check that the workspace is running on the selected runner
	if runner.ID != *workspace.RunnerID {
		utils.ErrorResponse(
			c,
			http.StatusNotFound,
			"workspace not found",
		)
		return
	}

	// start ssh client
	signer, err := ssh.ParsePrivateKey([]byte(workspace.User.SshPrivateKey))
	if err != nil {
		// TODO: log error
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connect to git server
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", sshHost), config)
	if err != nil {
		utils.ErrorResponse(
			c,
			http.StatusTeapot,
			"can't connect to the Git server",
		)
		return
	}
	defer client.Close()

	// start a new ssh session
	session, err := client.NewSession()
	if err != nil {
		// TODO: log error
		utils.ErrorResponse(
			c,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}
	defer session.Close()

	// upgrate inbound connection
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "cannot upgrade ws connection")
		return
	}
	defer wsConn.Close()

	stdinPipe, _ := session.StdinPipe()
	stdoutPipe, _ := session.StdoutPipe()

	cmd := fmt.Sprintf("%s %s", sshCmd, repoPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// read from ws and forward to stdin
	go func() {
		defer stdinPipe.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				wsConn.SetReadDeadline(time.Now().Add(2 * GitSSHPingInterval * time.Second))
				mt, data, err := wsConn.ReadMessage()
				if err != nil {
					cancel()
					return
				}

				if mt == websocket.BinaryMessage {
					if _, err := stdinPipe.Write(data); err != nil {
						cancel()
						return
					}
				}
			}
		}
	}()

	// read from stdout and forward to ws
	go func() {
		buf := make([]byte, 1024)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := stdoutPipe.Read(buf)
				if err != nil {
					cancel()
					return
				}

				if n > 0 {
					if err := wsConn.WriteMessage(
						websocket.BinaryMessage,
						buf[:n],
					); err != nil {
						cancel()
						return
					}
				}
			}
		}
	}()

	// ping ws connection to detect disconnection
	go func() {
		ticker := time.NewTicker(GitSSHPingInterval * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := wsConn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					cancel()
					return
				}
			}
		}
	}()

	// start the ssh command
	if err := session.Start(cmd); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to start SSH session")
		return
	}

	session.Wait()
	<-ctx.Done()

	// Add a small delay to allow goroutines to finish sending buffered data
	// before closing connections. This prevents connection drops during git push
	// confirmation messages.
	time.Sleep(GitSSHCloseDelay)

	wsConn.Close()
	session.Close()
}
