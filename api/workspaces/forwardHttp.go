package workspaces

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/davidebianchi03/codebox/api/utils"
	"github.com/davidebianchi03/codebox/db"
	"github.com/davidebianchi03/codebox/db/models"
	runnerinterface "github.com/davidebianchi03/codebox/runner-interface"
	"github.com/gin-gonic/gin"
)

func HandleForwardHttp(ctx *gin.Context) {
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
	containerName, found := ctx.Params.Get("containerName")
	if !found {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "workspace not found",
		})
		return
	}

	var workspace models.Workspace
	result := db.DB.Where(map[string]interface{}{"ID": workspaceId}).Preload("Runner").Find(&workspace)
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
	container := models.WorkspaceContainer{}
	r := db.DB.First(&container, map[string]interface{}{
		"workspace_id":   workspace.ID,
		"container_name": containerName,
	})

	if r.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	if container.ID <= 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "container not found, check that workspace is running and that you can connect to this container",
		})
		return
	}

	var port models.WorkspaceContainerPort
	db.DB.First(&port, map[string]interface{}{
		"container_id": container.ID,
		"port_number":  portNumber,
	})

	if port.ID <= 0 {
		// TODO: redirect to error page
		ctx.JSON(http.StatusNotFound, gin.H{
			"detail": "requested port is not forwarded",
		})
		return
	}

	if !port.Public {
		user, err := utils.GetUserFromContext(ctx)
		if err != nil {
			// TODO: redirect to login page
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"detail": "missing or invalid authorization token",
			})
			return
		}

		if user.ID != workspace.UserID {
			// TODO: redirect to login page
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"detail": "missing or invalid authorization token",
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

	parsedUrl, err := url.Parse(ctx.Request.URL.Path)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"detail": "internal server error",
		})
		return
	}

	ri := runnerinterface.RunnerInterface{
		Runner: workspace.Runner,
	}
	if err = ri.ForwardHttp(&workspace, &container, &port, parsedUrl.Query().Get("request_path"), ctx.Writer, ctx.Request); err != nil {
		ctx.JSON(http.StatusTeapot, gin.H{
			"detail": "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusTeapot, gin.H{
		"detail": "internal server error",
	})
}
