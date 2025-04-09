package wsv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	// chat := routerGroup.Group("/chat")
	// {
	// 	chat.GET("/room/:roomID/token/:token", app.CustomerControllers.ChatController.HandleWebsocket)
	// }
	routerGroup.GET("/chat/room/:roomID/token/:token")
	routerGroup.GET("/notifications/token/:token")
}
