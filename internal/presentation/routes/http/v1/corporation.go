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

	chat := routerGroup.Group("/chat")
	{
		chat.GET("/room/:corporationID", app.Controllers.Corporation.ChatController.GetRoom)
		chat.GET("/rooms/:corporationID", app.Controllers.Corporation.ChatController.GetRooms)
		chat.PUT("/room/:roomID/block", app.Controllers.Corporation.ChatController.BlockRoom)
		chat.PUT("/room/:roomID/unblock", app.Controllers.Corporation.ChatController.UnBlockRoom)
	}

	requests := routerGroup.Group(":corporationID/requests")
	{
		requests.GET("/installation", app.Controllers.Corporation.InstallationController.GetInstallationRequests)
	}

	panels := routerGroup.Group(":corporationID/panels")
	{
		panels.POST("add", app.Controllers.Corporation.InstallationController.AddPanel)
		panels.GET("list", app.Controllers.Corporation.InstallationController.GetCorporationPanels)
	}

	maintenance := routerGroup.Group(":corporationID/maintenance")
	{
		requests := maintenance.Group("/request")
		{
			requests.GET("/list", app.Controllers.Corporation.MaintenanceController.GetMaintenanceRequests)
			requests.POST("/handle", app.Controllers.Corporation.MaintenanceController.HandleMaintenanceRequest)
		}
		records := maintenance.Group("/record")
		{
			records.POST("/add", app.Controllers.Corporation.MaintenanceController.AddMaintenanceRecord)
			records.GET("/list", app.Controllers.Corporation.MaintenanceController.GetCorporationMaintenanceRecords)
			records.GET("/list/:panelID", app.Controllers.Corporation.MaintenanceController.GetCorporationMaintenanceRecordsByPanel)
		}
	}
}
