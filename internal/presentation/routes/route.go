package routes

import (
	httpv1 "github.com/BargheNo/Backend/internal/presentation/routes/http/v1"
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func Run(ginEngine *gin.Engine, app *wire.Application) {
	ginEngine.Use(app.Middlewares.Logger.GinLoggerMiddleware)
	ginEngine.Use(app.Middlewares.Localization.Localization)
	ginEngine.Use(app.Middlewares.Recovery.Recovery)
	ginEngine.Use(app.Middlewares.RateLimit.RateLimit)

	v1 := ginEngine.Group("/v1")
	registerSampleRoutes(v1, app)
}

func registerSampleRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	httpv1.SetupSampleRoutes(v1, app)
}
