package workspaces

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/davidebianchi03/codebox/api/utils"
	"github.com/davidebianchi03/codebox/config"
	dbconn "github.com/davidebianchi03/codebox/db/connection"
	"github.com/davidebianchi03/codebox/db/models"
	"github.com/davidebianchi03/codebox/runnerinterface"
	"github.com/gin-gonic/gin"
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
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
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
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "container not found, check that workspace is running and that you can connect to this container",
		})
		return
	}

	var port *models.WorkspaceContainerPort
	dbconn.DB.First(&port, map[string]interface{}{
		"container_id": container.ID,
		"port_number":  portNumber,
	})

	if port == nil {
		// TODO: redirect to error page
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "requested port is not forwarded",
		})
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
			// TODO: error 404
			ctx.JSON(http.StatusNotFound, gin.H{
				"detail": "not found",
			})
			return
		}
	}

	if workspace.Runner == nil {
		ctx.JSON(http.StatusTeapot, gin.H{
			"detail": "runner not found",
		})
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
