package middleware

import "github.com/gin-gonic/gin"

func CORSMiddleware(c *gin.Context) {
	httpOrigin := c.Request.Header["Origin"]
	if len(httpOrigin) > 0 {
		c.Writer.Header().Set("Access-Control-Allow-Origin", httpOrigin[0])
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	}

	c.Next()
}
