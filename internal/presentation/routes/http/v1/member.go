package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupMemberRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	corp := routerGroup.Group("/corp")
	{
		corp.POST("/change-password", app.Controllers.Member.CorporationController.ChangePassword)
		corp.POST("/update_contact_info", app.Controllers.Member.CorporationController.UpdateContactInfo)
		corp.GET("/installation_requests", app.Controllers.Member.CorporationController.GetInstallationRequests)
		corp.GET("/info", app.Controllers.Member.CorporationController.GetCorporationInfo)

		address := corp.Group("address")
		{
			address.POST("/add", app.Controllers.Member.CorporationController.AddAddress)
			address.POST("/edit", app.Controllers.Member.CorporationController.EditAddress)
			address.POST("/delete", app.Controllers.Member.CorporationController.DeleteAddress)
		}

		bid := corp.Group("bid")
		{
			bid.POST("/set", app.Controllers.Member.CorporationController.SetBid)
			bid.POST("/cancel", app.Controllers.Member.CorporationController.CancelBid)
			bid.GET("/get", app.Controllers.Member.CorporationController.GetBids)
		}
	}
}
