package middleware

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/davidebianchi03/codebox/api/workspaces"
	"github.com/davidebianchi03/codebox/config"
	"github.com/gin-gonic/gin"
)

func PortForwardingMiddleware(ctx *gin.Context) {
	// check if request hostname is a subdomain
	// to be valid a subdomain must start with codebox-- and must match the following
	// format codebox--<workspace_id>--<container_name>--<port_number>
	if strings.Index(ctx.Request.Host, fmt.Sprintf(".%s", config.Environment.ExternalUrl)) > 0 && strings.Index(ctx.Request.Host, "codebox--") == 0 {
		domainParts := strings.Split(ctx.Request.Host, ".")
		splittedSubDomain := strings.Split(domainParts[0], "--")

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
			ctx.Request.URL.Path = fmt.Sprintf(
				"/api/v1/workspace/%s/container/%s/forward-http/%s?request_path=%s",
				workspaceId,
				containerName,
				portNumber,
				url.QueryEscape(ctx.Request.URL.Path),
			)
		} else {
			ctx.Request.URL.Path = fmt.Sprintf(
				"/api/v1/workspace/%s/container/%s/forward-http/%s?request_path=%s",
				workspaceId,
				containerName,
				portNumber,
				url.QueryEscape(ctx.Request.URL.Path+"?"+ctx.Request.URL.RawQuery),
			)
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

		workspaces.HandleForwardHttp(ctx)
		ctx.Abort()
		return
	}

	ctx.Next()
}
