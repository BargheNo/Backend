package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupMemberRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	corp := routerGroup.Group("/corp")
	{
		corp.POST("/update_contact_info", app.Controllers.Member.CorporationController.UpdateContactInfo)
		corp.POST("/add_address", app.Controllers.Member.CorporationController.AddAddress)
		corp.POST("/edit_address", app.Controllers.Member.CorporationController.EditAddress)
		corp.POST("/delete_address", app.Controllers.Member.CorporationController.DeleteAddress)
		corp.GET("/installation_requests", app.Controllers.Member.CorporationController.GetInstallationRequests)
		corp.POST("/bid", app.Controllers.Member.CorporationController.SetBid)
		corp.POST("cancel_bid", app.Controllers.Member.CorporationController.CancelBid)
		corp.GET("/bids", app.Controllers.Member.CorporationController.GetBids)
	}
}
