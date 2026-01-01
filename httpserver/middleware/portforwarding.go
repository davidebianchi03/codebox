package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
	httperrors "gitlab.com/codebox4073715/codebox/httpserver/errors"
	"gitlab.com/codebox4073715/codebox/httpserver/ports"
)

// check if request hostname is a subdomain
// to be valid a subdomain must start with codebox-- and must match the following
// format codebox--<workspace_id>--<container_name>--<port_number>
func PortForwardingMiddleware(ctx *gin.Context) {
	requestDomain := ctx.Request.Host
	if strings.Contains(requestDomain, fmt.Sprintf(".%s", config.Environment.WildcardDomain)) {
		subdomains := strings.Split(strings.ReplaceAll(requestDomain, fmt.Sprintf(".%s", config.Environment.WildcardDomain), ""), ".")
		if len(subdomains) == 0 {
			httperrors.RenderError(ctx, http.StatusNotFound, "Not found")
			ctx.Abort()
			return
		}

		portSubdomain := subdomains[len(subdomains)-1]
		if !strings.HasPrefix(portSubdomain, "codebox--") {
			httperrors.RenderError(ctx, http.StatusNotFound, "Not found")
			ctx.Abort()
			return
		}

		splittedSubDomain := strings.Split(portSubdomain, "--")

		if len(splittedSubDomain) != 4 {
			httperrors.RenderError(ctx, http.StatusNotFound, "Not found")
			ctx.Abort()
			return
		}

		if ctx.Request.URL.Path == fmt.Sprintf(
			"/api/v1/auth/subdomains/callback-%s",
			url.PathEscape(config.Environment.AuthCookieName),
		) {
			// TODO: check if workspace exists
			ctx.Next()
			return
		}

		workspaceId, err := strconv.Atoi(splittedSubDomain[1])
		if err != nil || workspaceId <= 0 {
			httperrors.RenderError(ctx, http.StatusNotFound, "Not found")
			ctx.Abort()
			return
		}

		portNumber, err := strconv.Atoi(splittedSubDomain[3])
		if err != nil || portNumber <= 0 || portNumber >= 65536 {
			httperrors.RenderError(ctx, http.StatusNotFound, "Not found")
			ctx.Abort()
			return
		}

		ports.ForwardHttpPort(
			ctx,
			uint(workspaceId),
			splittedSubDomain[2],
			uint(portNumber),
			ctx.Request.URL.String(),
		)

		ctx.Abort()
		return
	}

	ctx.Next()
}
