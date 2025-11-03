package router

import (
	"path"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/api"
	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/views"
)

func SetupRouter() *gin.Engine {
	if config.Environment.DebugEnabled {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	api.V1ApiRoutes(r)
	views.ViewsRoutes(r)
	r.LoadHTMLGlob(path.Join(config.Environment.BaseDir, "html", "templates", "*"))
	return r
}
