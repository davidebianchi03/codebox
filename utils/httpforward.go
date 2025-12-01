package utils

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
)

/*
Handle Http port forward
*/
func ForwardHttpPort(
	ctx *gin.Context,
	workspaceId uint,
	containerName string,
	portNumber uint,
	path string,
) {
	// retrieve workspace details
	workspace, err := models.RetrieveWorkspaceById(workspaceId)
	if err != nil {
		RenderError(ctx, http.StatusInternalServerError, "Unknown error")
		return
	}

	if workspace == nil {
		RenderError(ctx, http.StatusNotFound, "Workspace not found")
		return
	}

	// retrieve container
	container, err := models.RetrieveWorkspaceContainerByName(
		*workspace,
		containerName,
	)

	if err != nil {
		RenderError(ctx, http.StatusInternalServerError, "Unknown error")
		return
	}

	if container == nil {
		RenderError(ctx, http.StatusNotFound, "Container not found")
		return
	}

	// retrieve port
	port, err := models.RetrieveContainerPortByPortNumber(
		*container,
		portNumber,
	)
	if err != nil {
		RenderError(ctx, http.StatusInternalServerError, "Unknown error")
		return
	}

	if port == nil {
		RenderError(
			ctx,
			http.StatusNotFound,
			"Port is not exposed or you don't have the permission to view it",
		)
		return
	}

	if !port.Public {
		user, err := utils.GetUserFromContext(ctx)
		if err != nil {
			requestProtocol := ctx.Request.Header.Get("X-Forwarded-Proto")
			if requestProtocol != "http" && requestProtocol != "https" {
				requestProtocol = "http"
			}
			currentLocation := fmt.Sprintf("%s://%s%s", requestProtocol, ctx.Request.Host, ctx.Request.URL.String())
			authorizeEndpointUrl := ""

			if config.Environment.UseSubDomains {
				authorizeEndpointUrl = fmt.Sprintf(
					"%s/api/v1/auth/subdomains/authorize?next=%s",
					config.Environment.ExternalUrl,
					url.QueryEscape(currentLocation),
				)
			} else {
				authorizeEndpointUrl = fmt.Sprintf("/login?next=%s", url.QueryEscape(currentLocation))
			}

			ctx.Redirect(
				http.StatusTemporaryRedirect,
				authorizeEndpointUrl,
			)
			return
		}

		if user.ID != workspace.UserID {
			RenderError(ctx, http.StatusNotFound, "Port is not exposed or you don't have the permission to view it")
			return
		}
	}

	if workspace.Runner == nil {
		RenderError(ctx, http.StatusInternalServerError, "Unknown error")
		return
	}

	ri := runnerinterface.RunnerInterface{
		Runner: workspace.Runner,
	}

	if err := ri.ForwardHttp(workspace, container, port, path, ctx.Writer, ctx.Request); err != nil {
		RenderError(ctx, http.StatusInternalServerError, "Unknown error")
		return
	}
}
