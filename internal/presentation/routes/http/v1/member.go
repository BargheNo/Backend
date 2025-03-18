package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupMemberRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	corp := routerGroup.Group("/corp")
	{
		corp.POST("/add_contact_info", app.Controllers.Member.CorporationController.AddContactInfo)
		corp.GET("/installation_requests", app.Controllers.Member.CorporationController.GetInstallationRequests)
		corp.POST("/bid", app.Controllers.Member.CorporationController.SetBid)
		corp.POST("cancel_bid", app.Controllers.Member.CorporationController.CancelBid)
		corp.GET("/bids", app.Controllers.Member.CorporationController.GetBids)
	}
}
