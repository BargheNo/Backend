package wsv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	routerGroup.GET("/chat/room/:roomID/token/:token", app.Controllers.Customer.ChatController.HandleWebsocket)
	routerGroup.GET("/notifications/token/:token", app.Controllers.Customer.NotificationController.HandleWebsocket)
	routerGroup.GET("/monitoring/panel/:panelID/token/:token", app.Controllers.Customer.MonitoringController.HandleWebsocket)
}
