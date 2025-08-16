package routes

import (
	httpv1 "github.com/BargheNo/Backend/internal/presentation/routes/http/v1"
	wsv1 "github.com/BargheNo/Backend/internal/presentation/routes/ws/v1"
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(ginEngine *gin.Engine, app *wire.Application) {
	ginEngine.Use(app.Middlewares.CORS.CORS())
	ginEngine.Use(app.Middlewares.Logger.GinLoggerMiddleware)
	ginEngine.Use(app.Middlewares.Localization.Localization)
	ginEngine.Use(app.Middlewares.Recovery.Recovery)
	ginEngine.Use(app.Middlewares.RateLimit.RateLimit)
	ginEngine.Use(app.Middlewares.Prometheus.PrometheusMiddleware)

	ginEngine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := ginEngine.Group("/v1")
	registerGeneralRoutes(v1, app)
	registerCustomerRoutes(v1, app)
	registerCorporationRoutes(v1, app)
	registerAdminRoutes(v1, app)
}

func registerGeneralRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	httpv1.SetupGeneralRoutes(v1, app)
}

func registerCustomerRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	user := v1.Group("/user")
	user.Use(app.Middlewares.Authentication.AuthRequired)
	httpv1.SetupCustomerRoutes(user, app)

	wsUser := v1.Group("/user")
	wsUser.Use(app.Middlewares.WebsocketMiddleware.UpgradeToWebSocket)
	wsv1.SetupCustomerRoutes(wsUser, app)
}

func registerCorporationRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	corporation := v1.Group("/corp/:corporationID")
	corporation.Use(app.Middlewares.Authentication.AuthRequired)
	corporation.Use(app.Middlewares.Authentication.CorporationAccessRequired)
	httpv1.SetupCorporationRoutes(corporation, app)
}

func registerAdminRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	admin := v1.Group("/admin")
	admin.Use(app.Middlewares.Authentication.AuthRequired)
	httpv1.SetupAdminRoutes(admin, app)
}
