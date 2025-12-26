package views

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	httperrors "gitlab.com/codebox4073715/codebox/httpserver/errors"
	"gitlab.com/codebox4073715/codebox/httpserver/ports"
)

/*
View that handles port forward
Forward HTTP ports
*/
func HandlePortForwardView(c *gin.Context) {
	workspaceIdStr, found := c.Params.Get("workspaceId")
	if !found {
		httperrors.RenderError(c, http.StatusInternalServerError, "Unknown error")
		return
	}

	workspaceId, err := strconv.Atoi(workspaceIdStr)
	if err != nil || workspaceId <= 0 {
		httperrors.RenderError(c, http.StatusNotFound, "Workspace not found")
	}

	containerName, found := c.Params.Get("containerName")
	if !found {
		httperrors.RenderError(c, http.StatusInternalServerError, "Unknown error")
		return
	}

	portNumberStr, found := c.Params.Get("portNumber")
	if !found {
		httperrors.RenderError(c, http.StatusInternalServerError, "Unknown error")
		return
	}

	portNumber, err := strconv.Atoi(portNumberStr)
	if err != nil || portNumber <= 0 || portNumber >= 65536 {
		httperrors.RenderError(c, http.StatusBadRequest, "Invalid port number")
	}

	url := c.Request.URL.Path
	if len(c.Request.URL.RawQuery) > 0 {
		url += "?" + c.Request.URL.RawQuery
	}

	ports.ForwardHttpPort(
		c,
		uint(workspaceId),
		containerName,
		uint(portNumber),
		c.Request.URL.Path,
	)
}
