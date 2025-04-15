package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	ticket := routerGroup.Group("/ticket")
	{
		ticket.GET("", app.Controllers.Admin.TicketController.GetTickets)
		ticket.GET("/:ticketID/comments", app.Controllers.Admin.TicketController.GetComments)
	}
}
