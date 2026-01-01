package httpserver

import (
	"path"

	"github.com/gin-gonic/gin"
	"gitlab.com/codebox4073715/codebox/config"
	"gitlab.com/codebox4073715/codebox/httpserver/api"
	"gitlab.com/codebox4073715/codebox/httpserver/middleware"
	"gitlab.com/codebox4073715/codebox/httpserver/views"
)

func SetupRouter() *gin.Engine {
	if config.Environment.DebugEnabled {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// middlewares
	r.Use(middleware.PortForwardingMiddleware)
	r.Use(middleware.CORSMiddleware)

	// apis
	api.V1ApiRoutes(r)

	// views
	views.ViewsRoutes(r)
	r.LoadHTMLGlob(
		path.Join(config.Environment.TemplatesFolder, "templates", "*"),
	)
	return r
}
