package cli

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleDownloadCLI(ctx *gin.Context) {
	os := ctx.Query("os")
	if os == "windows" {
		ctx.File("./codebox-cli-windows-amd64.exe")
		ctx.Header("Content-Disposition", "attachment; filename=codebox-cli-windows-amd64.exe")
		ctx.Header("Content-Type", "application/gzip")
		ctx.Status(http.StatusOK)
	} else if os == "linux" {
		ctx.File("./codebox-cli-linux-amd64.exe")
		ctx.Header("Content-Disposition", "attachment; filename=codebox-cli-linux-amd64")
		ctx.Header("Content-Type", "application/gzip")
		ctx.Status(http.StatusOK)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"detail": "invalid or missing url query argument 'os'",
		})
	}
}
