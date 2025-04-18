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
		ticket.POST("/:ticketID/comments", app.Controllers.Admin.TicketController.CreateComment)
		ticket.POST("/:ticketID/resolve", app.Controllers.Admin.TicketController.ResolveTicket)
	}

	report := routerGroup.Group("/report")
	{
		report.GET("/maintenance", app.Controllers.Admin.ReportController.GetMaintenanceReports)
		report.GET("/panel", app.Controllers.Admin.ReportController.GetPanelReports)
		report.POST("/resolve/:reportID", app.Controllers.Admin.ReportController.ResolveReport)
	}
}
