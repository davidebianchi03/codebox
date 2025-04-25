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

	requestDomain := ctx.Request.Host
	// if strings.Contains(requestDomain, fmt.Sprintf(".%s", config.Environment.WildcardExternalUrl)) {
	if strings.Contains(requestDomain, fmt.Sprintf(config.Environment.WildcardExternalUrl)) { // TODO: use previous line
		requestDomain = "codebox--2--phpmyadmin--80.codebox--2--codebox--8080.codebox.davidebianchi.eu" // TODO: remove
		subdomains := strings.Split(strings.ReplaceAll(requestDomain, fmt.Sprintf(".%s", config.Environment.WildcardExternalUrl), ""), ".")
		if len(subdomains) == 0 {
			// TODO: show error page: 404
			return
		}

		portSubdomain := subdomains[len(subdomains)-1]
		if !strings.HasPrefix(portSubdomain, "codebox--") {
			// TODO: show error page: 404
			return
		}

		splittedSubDomain := strings.Split(portSubdomain, "--")

		if len(splittedSubDomain) != 4 {
			ctx.JSON(400, gin.H{
				"detail": "invalid hostname",
			})
			ctx.Abort()
			return
		}

		if ctx.Request.URL.Path == fmt.Sprintf("/api/v1/auth/subdomains/callback-%s", url.PathEscape(config.Environment.AuthCookieName)) {
			// TODO: check if workspace exists
			ctx.Next()
			return
		}

		workspaces.ForwardHttpPort(
			ctx,
			splittedSubDomain[1],
			splittedSubDomain[2],
			splittedSubDomain[3],
			ctx.Request.URL.String(),
		)

		ctx.Abort()
		return
	}

	ctx.Next()
}
