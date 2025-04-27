package workspaces

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	apierrors "gitlab.com/codebox4073715/codebox/api/errors"
	"gitlab.com/codebox4073715/codebox/api/utils"
	"gitlab.com/codebox4073715/codebox/config"
	dbconn "gitlab.com/codebox4073715/codebox/db/connection"
	"gitlab.com/codebox4073715/codebox/db/models"
	"gitlab.com/codebox4073715/codebox/runnerinterface"
)

func HandleForwardHttp(ctx *gin.Context) {
	portNumber, found := ctx.Params.Get("portNumber")
	if !found {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "missing arg portNumber",
		})
		return
	}

	workspaceId, found := ctx.Params.Get("workspaceId")
	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "missing arg workspaceId",
		})
		return
	}

	containerName, found := ctx.Params.Get("containerName")
	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "missing arg containerName",
		})
		return
	}

	parsedUrl, err := url.Parse(ctx.Request.URL.Path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	ForwardHttpPort(ctx, workspaceId, containerName, portNumber, parsedUrl.Query().Get("request_path"))
}

func ForwardHttpPort(ctx *gin.Context, workspaceId string, containerName string, portNumber string, path string) {
	var workspace *models.Workspace
	result := dbconn.DB.Where(map[string]interface{}{"ID": workspaceId}).Preload("Runner").Find(&workspace)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if workspace == nil {
		apierrors.RenderError(
			ctx, http.StatusNotFound, "Workspace does not exists or you don't have the permission to view it",
		)
		return
	}

	// retrieve development container
	var container *models.WorkspaceContainer
	r := dbconn.DB.Find(&container, map[string]interface{}{
		"workspace_id":   workspace.ID,
		"container_name": containerName,
	})

	if r.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if container == nil {
		apierrors.RenderError(
			ctx, http.StatusNotFound, "Workspace does not exists or you don't have the permission to view it",
		)
		return
	}

	var port *models.WorkspaceContainerPort
	dbconn.DB.First(&port, map[string]interface{}{
		"container_id": container.ID,
		"port_number":  portNumber,
	})

	if port == nil {
		apierrors.RenderError(
			ctx, http.StatusNotFound, "Port is not exposed or you don't have the permission to view it",
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
			authorizeEndpointUrl := fmt.Sprintf("%s/api/v1/auth/subdomains/authorize?next=%s", config.Environment.ExternalUrl, url.QueryEscape(currentLocation))
			ctx.Redirect(
				http.StatusTemporaryRedirect,
				authorizeEndpointUrl,
			)
			return
		}

		if user.ID != workspace.UserID {
			apierrors.RenderError(
				ctx, http.StatusNotFound, "Port is not exposed or you don't have the permission to view it",
			)
			return
		}
	}

	if workspace.Runner == nil {
		apierrors.RenderError(
			ctx, http.StatusNotFound, "Runner not found",
		)
		return
	}

	ri := runnerinterface.RunnerInterface{
		Runner: workspace.Runner,
	}

	if err := ri.ForwardHttp(workspace, container, port, path, ctx.Writer, ctx.Request); err != nil {
		ctx.JSON(http.StatusTeapot, gin.H{
			"detail": "internal server error",
		})
		return
	}
}
