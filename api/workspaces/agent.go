package workspaces

import (
	"fmt"
	"net/http"

	"codebox.com/db"
	"codebox.com/proxy"
	"github.com/gin-gonic/gin"
)

func HandleForwardContainerPort(ctx *gin.Context) {
	// TODO: autenticazione se la porta non è pubblica
	// TODO: security, non permettere la connessione a porte nn consentite
	// user, err := utils.GetUserFromContext(ctx)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"detail": "internal server error",
	// 	})
	// 	return
	// }

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

	proxyTargetUrl := fmt.Sprintf("http://%s:%d", developmentContainer.ExternalIPv4, developmentContainer.AgentExternalPort)

	proxyHeaders := map[string][]string{}

	proxy, err := proxy.CreateReverseProxy(proxyTargetUrl, 30, 30, true, proxyHeaders)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}
