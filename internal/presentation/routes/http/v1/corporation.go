package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCorporationRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	// TODO: fix API URLs
	auth := routerGroup.Group("/auth") // cancel
	{
		auth.POST("/corporation/reset-password", app.Controllers.Corporation.CorporationController.ChangePassword)
	}

	addresses := routerGroup.Group("/address")
	{
		addresses.POST("/corp", app.Controllers.Corporation.AddressController.CreateCorporationAddress)
		addresses.DELETE("/corp", app.Controllers.Corporation.AddressController.DeleteCorporationAddress)
	}

	bids := routerGroup.Group("/bids")
	{
		bids.POST("/set", app.Controllers.Corporation.BidController.SetBid)
		bids.PUT("/cancel", app.Controllers.Corporation.BidController.CancelBid)
		bids.GET("/list", app.Controllers.Corporation.BidController.GetBids)
	}

	corp := routerGroup.Group("/corp")
	{
		corp.GET("/info", app.Controllers.Corporation.CorporationController.GetCorporationInfo)
		corp.POST("/contact-info", app.Controllers.Corporation.CorporationController.UpdateContactInfo)
		corp.GET("/installation", app.Controllers.Corporation.InstallationController.GetInstallationRequests)
	}
}
