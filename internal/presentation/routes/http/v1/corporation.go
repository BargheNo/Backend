package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCorporationRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	profile := routerGroup.Group("/:corporationID/profile")
	{
		profile.GET("", app.Controllers.Corporation.CorporationController.GetMyProfile)
		profile.POST("/address", app.Controllers.Corporation.CorporationController.AddAddress)
		profile.DELETE("/address/:addressID", app.Controllers.Corporation.CorporationController.DeleteAddress)
		profile.POST("/contacts", app.Controllers.Corporation.CorporationController.AddContactInformation)
		profile.DELETE("/contacts/:contactID", app.Controllers.Corporation.CorporationController.DeleteContactInformation)
		profile.PUT("/logo", app.Controllers.Corporation.CorporationController.ChangeLogo)
	}

	guarantees := routerGroup.Group("/:corporationID/guarantee")
	{
		guarantees.GET("", app.Controllers.Corporation.GuaranteeController.GetGuarantees)
		guarantees.GET("/type", app.Controllers.Corporation.GuaranteeController.GetGuaranteeTypes)
		guarantees.POST("", app.Controllers.Corporation.GuaranteeController.CreateGuarantee)
		guaranteesSubGroup := guarantees.Group("/:guaranteeID")
		{
			guaranteesSubGroup.GET("", app.Controllers.Corporation.GuaranteeController.GetGuarantee)
			guaranteesSubGroup.PUT("/status", app.Controllers.Corporation.GuaranteeController.UpdateGuarantee)
		}
	}

	installations := routerGroup.Group("/:corporationID/installation")
	{
		requests := installations.Group("/request")
		{
			requests.GET("", app.Controllers.Corporation.InstallationController.GetInstallationRequests)
			requestSubGroup := requests.Group("/:requestID")
			{
				requestSubGroup.GET("", app.Controllers.Corporation.InstallationController.GetInstallationRequest)
				requestSubGroup.POST("/bid", app.Controllers.Corporation.BidController.SetBid)
			}
		}
		panels := installations.Group("/panel")
		{
			panels.POST("", app.Controllers.Corporation.InstallationController.AddPanel)
			panels.GET("", app.Controllers.Corporation.InstallationController.GetCorporationPanels)
			panelsSubGroup := panels.Group("/:panelID")
			{
				panelsSubGroup.GET("", app.Controllers.Corporation.InstallationController.GetCorporationPanel)
				panelsSubGroup.PUT("/complete", app.Controllers.Corporation.InstallationController.CompleteInstallation)
			}
		}
	}

	bids := routerGroup.Group(":corporationID/bid")
	{
		bids.GET("", app.Controllers.Corporation.BidController.GetBids)
		bidsSubGroup := bids.Group("/:bidID")
		{
			bidsSubGroup.GET("", app.Controllers.Corporation.BidController.GetBid)
			bidsSubGroup.PUT("", app.Controllers.Corporation.BidController.UpdateBid)
			bidsSubGroup.PUT("/cancel", app.Controllers.Corporation.BidController.CancelBid)
		}
	}

	chat := routerGroup.Group("/chat")
	{
		chat.GET("/room/:corporationID", app.Controllers.Corporation.ChatController.GetRoom)
		chat.GET("/rooms/:corporationID", app.Controllers.Corporation.ChatController.GetRooms)
		chat.PUT("/room/:roomID/block", app.Controllers.Corporation.ChatController.BlockRoom)
		chat.PUT("/room/:roomID/unblock", app.Controllers.Corporation.ChatController.UnBlockRoom)
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
