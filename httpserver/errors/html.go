package utils

import (
	"github.com/gin-gonic/gin"
)

/*
Render errors with html templates
*/
func RenderError(c *gin.Context, status int, message string) {
	c.HTML(
		status,
		"errors.html",
		gin.H{
			"title":   "Oops! Something Went Wrong",
			"message": message,
		},
	)
}
