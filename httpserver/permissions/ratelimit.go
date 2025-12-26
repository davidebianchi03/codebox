package permissions

import "github.com/gin-gonic/gin"

func RateLimitedRoute(
	handler gin.HandlerFunc,
	callsPerPeriod int,
	periodSeconds int,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		ipAddress := c.ClientIP()
		requestPath := c.FullPath()

		_ = ipAddress
		_ = requestPath

		handler(c)
	}
}
