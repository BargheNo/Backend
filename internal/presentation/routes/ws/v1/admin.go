package wsv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	routerGroup.GET("/monitoring/panel/:panelID/token/:token", app.Controllers.Admin.MonitoringController.HandleWebsocket)
}
