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
	registerGeneralRoutes(v1, app)
	registerMemberRoutes(v1, app)
}

func registerGeneralRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	httpv1.SetupGeneralRoutes(v1, app)
}

func registerMemberRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	httpv1.SetupMemberRoutes(v1, app)
}
