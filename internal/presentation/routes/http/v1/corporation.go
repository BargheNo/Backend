package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCorporationRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	addresses := routerGroup.Group(":corporationID/address")
	{
		addresses.POST("", app.Controllers.Corporation.CorporationController.AddAddress)
		addresses.DELETE("", app.Controllers.Corporation.CorporationController.DeleteAddress)
	}

	profileComplete := routerGroup.Group(":corporationID")
	{
		profileComplete.POST("/certificates", app.Controllers.Corporation.CorporationController.SubmitCertificateFiles)
		profileComplete.POST("/contact", app.Controllers.Corporation.CorporationController.AddContactInformation)
	}

	bids := routerGroup.Group(":corporationID/bids")
	{
		bids.POST("/set", app.Controllers.Corporation.BidController.SetBid)
		bids.PUT("/cancel", app.Controllers.Corporation.BidController.CancelBid)
		bids.GET("/list", app.Controllers.Corporation.BidController.GetBids)
	}

	requests := routerGroup.Group("/requests")
	{
		requests.GET("/installation", app.Controllers.Corporation.InstallationController.GetInstallationRequests)
	}
}
