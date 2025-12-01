package views

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/utils"
)

/*
View that handles port forward
Forward HTTP ports
*/
func HandlePortForwardView(ctx *gin.Context) {
	workspaceIdStr, found := ctx.Params.Get("workspaceId")
	if !found {
		utils.RenderError(ctx, http.StatusInternalServerError, "Unknown error")
		return
	}

	workspaceId, err := strconv.Atoi(workspaceIdStr)
	if err != nil || workspaceId <= 0 {
		utils.RenderError(ctx, http.StatusNotFound, "Workspace not found")
	}

	containerName, found := ctx.Params.Get("containerName")
	if !found {
		utils.RenderError(ctx, http.StatusInternalServerError, "Unknown error")
		return
	}

	portNumberStr, found := ctx.Params.Get("portNumber")
	if !found {
		utils.RenderError(ctx, http.StatusInternalServerError, "Unknown error")
		return
	}

	portNumber, err := strconv.Atoi(portNumberStr)
	if err != nil || portNumber <= 0 || portNumber >= 65536 {
		utils.RenderError(ctx, http.StatusBadRequest, "Invalid port number")
	}

	url := ctx.Request.URL.Path
	if len(ctx.Request.URL.RawQuery) > 0 {
		url += "?" + ctx.Request.URL.RawQuery
	}

	utils.ForwardHttpPort(
		ctx,
		uint(workspaceId),
		containerName,
		uint(portNumber),
		ctx.Request.URL.Path,
	)
}
