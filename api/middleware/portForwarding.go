package middleware

import (
	"fmt"
	"net/url"
	"strings"

	"codebox.com/api/workspaces"
	"github.com/gin-gonic/gin"
)

func PortForwardingMiddleware(ctx *gin.Context) {
	ctx.Request.Host = "codebox--3--phpmyadmin--80.codebox.davidebianchi.eu"

	if len(strings.Split(ctx.Request.Host, ".")) > 1 {
		codeboxSubDomain := strings.Split(ctx.Request.Host, ".")[0]

		if strings.HasPrefix(codeboxSubDomain, "codebox--") {
			splittedSubDomain := strings.Split(codeboxSubDomain, "--")

			if len(splittedSubDomain) != 4 {
				ctx.JSON(400, gin.H{
					"detail": "invalid hostname",
				})
				ctx.Abort()
				return
			}

			workspaceId := splittedSubDomain[1]
			containerName := splittedSubDomain[2]
			portNumber := splittedSubDomain[3]

			if ctx.Request.URL.RawQuery == "" {
				ctx.Request.URL.Path = fmt.Sprintf("/api/v1/workspace/%s/container/%s/forward/%s?request_path=%s", workspaceId, containerName, portNumber, url.QueryEscape(ctx.Request.URL.Path))
			} else {
				ctx.Request.URL.Path = fmt.Sprintf("/api/v1/workspace/%s/container/%s/forward/%s?request_path=%s", workspaceId, containerName, portNumber, url.QueryEscape(ctx.Request.URL.Path+"?"+ctx.Request.URL.RawQuery))
			}

			newRequestParams := []gin.Param{
				{
					Key:   "workspaceId",
					Value: workspaceId,
				},
				{
					Key:   "containerName",
					Value: containerName,
				},
				{
					Key:   "portNumber",
					Value: portNumber,
				},
			}
			ctx.Params = newRequestParams

			workspaces.HandleForwardContainerPort(ctx)
			ctx.Abort()
			return
		}
	}

	ctx.Next()
}
