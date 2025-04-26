package errors

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// render errors
// if request comes from a browser return pretty
// html otherwise return json response
func RenderError(c *gin.Context, status int, message string) {
	userAgent := c.Request.UserAgent()
	if strings.Contains(userAgent, "Chrome") ||
		strings.Contains(userAgent, "Firefox") ||
		strings.Contains(userAgent, "Safari") ||
		strings.Contains(userAgent, "AppleWebKit") ||
		strings.Contains(userAgent, "Mozilla") {
		c.HTML(
			status,
			"errors.html",
			gin.H{
				"title":   "Oops! Something Went Wrong",
				"message": message,
			},
		)
	} else {
		c.JSON(status, gin.H{
			"details": message,
		})
	}
}
