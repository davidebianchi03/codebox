package workspaces

import (
	"fmt"
	"net/http"
	"strconv"

	"codebox.com/api/utils"
	"codebox.com/db"
	"codebox.com/proxy"
	"github.com/gin-gonic/gin"
)

func HandleForwardContainerPort(ctx *gin.Context) {
	portNumberStr, found := ctx.Params.Get("portNumber")
	if !found {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid port number",
		})
		return
	}

	portNumber, err := strconv.Atoi(portNumberStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid port number",
		})
		return
	}

	workspaceId, found := ctx.Params.Get("workspaceId")
	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}
	containerId, found := ctx.Params.Get("containerId")
	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	var workspace db.Workspace
	result := db.DB.Where(map[string]interface{}{"ID": workspaceId /* "owner_id": user.ID*/}).Find(&workspace)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	// retrieve development container
	developmentContainer := db.WorkspaceContainer{}
	result = db.DB.Where(
		map[string]interface{}{"workspace_id": workspace.ID, "can_connect_remote_developing": true, "ID": containerId}).Preload("ForwardedPorts").Find(&developmentContainer)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if result.RowsAffected != 1 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "container not found, check that workspace is running and that you can connect to this container",
		})
		return
	}

	// check that requested port is forwarded
	portFound := false
	var forwardedPort db.ForwardedPort
	for _, port := range developmentContainer.ForwardedPorts {
		if port.PortNumber == uint(portNumber) {
			forwardedPort = port
			portFound = true
			break
		}
	}

	if !portFound {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "port not forwarded",
		})
		return
	}

	if !forwardedPort.Active {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "port not active",
		})
		return
	}

	if !forwardedPort.Public {
		user, err := utils.GetUserFromContext(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"detail": "missing or invalid authorization token",
			})
			return
		}

		if user.ID != workspace.OwnerId {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"detail": "missing or invalid authorization token",
			})
			return
		}
	}

	proxyTargetUrl := fmt.Sprintf("http://%s:%d", developmentContainer.ExternalIPv4, developmentContainer.AgentExternalPort)

	proxyHeaders := http.Header{}
	proxyHeaders.Set("X-CodeBox-Forward-Host", "127.0.0.1")
	proxyHeaders.Set("X-CodeBox-Forward-Port", strconv.Itoa(int(forwardedPort.PortNumber)))
	proxyHeaders.Set("X-CodeBox-Forward-Domain", "localhost")
	if forwardedPort.ConnectionType == db.ConnectionTypeWS {
		proxyHeaders.Set("X-CodeBox-Forward-Scheme", "tcp_over_websocket")
	} else {
		proxyHeaders.Set("X-CodeBox-Forward-Scheme", "http")
	}

	proxy, err := proxy.CreateReverseProxy(proxyTargetUrl, 30, 30, true, proxyHeaders)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
