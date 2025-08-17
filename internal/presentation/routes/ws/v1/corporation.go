package wsv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCorporationRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	routerGroup.GET("/:corporationID/monitoring/panel/:panelID/token/:token", app.Controllers.Corporation.MonitoringController.HandleWebsocket)
}
