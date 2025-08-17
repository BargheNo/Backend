package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCorporationRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	const status string = "/status"

	profile := routerGroup.Group("/profile")
	{
		profile.GET("", app.Controllers.Corporation.CorporationController.GetMyProfile)
		profile.GET("/public", app.Controllers.Corporation.CorporationController.GetPublicProfile)
		profile.POST("/address", app.Controllers.Corporation.CorporationController.AddAddress)
		profile.DELETE("/address/:addressID", app.Controllers.Corporation.CorporationController.DeleteAddress)
		profile.POST("/contacts", app.Controllers.Corporation.CorporationController.AddContactInformation)
		profile.DELETE("/contacts/:contactID", app.Controllers.Corporation.CorporationController.DeleteContactInformation)
		profile.PUT("/logo", app.Controllers.Corporation.CorporationController.ChangeLogo)
		profile.PUT("", app.Controllers.Corporation.CorporationController.UpdateRegister)
		profile.PUT("/certificates", app.Controllers.Corporation.CorporationController.SubmitCertificateFiles)
	}

	staff := routerGroup.Group("/staff")
	{
		staff.GET(status, app.Controllers.Corporation.CorporationController.GetStaffStatuses)
		staff.GET("", app.Controllers.Corporation.CorporationController.GetStaffList)
		staff.POST("", app.Controllers.Corporation.CorporationController.CreateCorporationStaff)
		staffSubGroup := staff.Group("/:staffID")
		{
			staffSubGroup.GET("", app.Controllers.Corporation.CorporationController.GetStaff)
			staffSubGroup.PUT("", app.Controllers.Corporation.CorporationController.EditCorporationStaff)
		}
	}

	guarantees := routerGroup.Group("/guarantee")
	{
		guarantees.GET("", app.Controllers.Corporation.GuaranteeController.GetGuarantees)
		guarantees.GET("/type", app.Controllers.Corporation.GuaranteeController.GetGuaranteeTypes)
		guarantees.POST("", app.Controllers.Corporation.GuaranteeController.CreateGuarantee)
		guaranteesSubGroup := guarantees.Group("/:guaranteeID")
		{
			guaranteesSubGroup.GET("", app.Controllers.Corporation.GuaranteeController.GetGuarantee)
			guaranteesSubGroup.PUT(status, app.Controllers.Corporation.GuaranteeController.UpdateGuarantee)
		}
	}

	installations := routerGroup.Group("/installation")
	{
		requests := installations.Group("/request")
		{
			requests.GET("", app.Controllers.Corporation.InstallationController.GetInstallationRequests)
			requests.GET(status, app.Controllers.Corporation.MaintenanceController.GetMaintenanceStatuses)
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
				guaranteeViolation := panelsSubGroup.Group("/guarantee/violation")
				{
					guaranteeViolation.POST("", app.Controllers.Corporation.InstallationController.ViolatePanelGuarantee)
					guaranteeViolation.GET("", app.Controllers.Corporation.InstallationController.GetPanelGuaranteeViolation)
					guaranteeViolation.DELETE("", app.Controllers.Corporation.InstallationController.ClearPanelGuaranteeViolation)
					guaranteeViolation.PUT("", app.Controllers.Corporation.InstallationController.UpdatePanelGuaranteeViolation)
				}
			}
		}
	}

	maintenances := routerGroup.Group("/maintenance")
	{
		requests := maintenances.Group("/request")
		{
			requests.GET(status, app.Controllers.Corporation.MaintenanceController.GetMaintenanceStatuses)
			requests.GET("", app.Controllers.Corporation.MaintenanceController.GetAllMaintenanceRequests)
			requestsSubGroup := requests.Group("/:requestID")
			{
				requestsSubGroup.GET("", app.Controllers.Corporation.MaintenanceController.GetMaintenanceRequest)
				requestsSubGroup.PUT("/accept", app.Controllers.Corporation.MaintenanceController.AcceptMaintenanceRequest)
				requestsSubGroup.PUT("/reject", app.Controllers.Corporation.MaintenanceController.RejectMaintenanceRequest)
				records := requestsSubGroup.Group("/record")
				{
					records.POST("", app.Controllers.Corporation.MaintenanceController.CreateMaintenanceRecord)
					records.PUT("", app.Controllers.Corporation.MaintenanceController.UpdateMaintenanceRecord)
				}
			}
		}
	}

	bids := routerGroup.Group("/bid")
	{
		bids.GET("", app.Controllers.Corporation.BidController.GetBids)
		bids.GET(status, app.Controllers.Corporation.BidController.GetBidStatuses)
		bidsSubGroup := bids.Group("/:bidID")
		{
			bidsSubGroup.GET("", app.Controllers.Corporation.BidController.GetBid)
			bidsSubGroup.PUT("", app.Controllers.Corporation.BidController.UpdateBid)
			bidsSubGroup.PUT("/cancel", app.Controllers.Corporation.BidController.CancelBid)
		}
	}

	chat := routerGroup.Group("/chat")
	{
		chat.GET("/room", app.Controllers.Corporation.ChatController.GetRoom)
		chat.GET("/rooms", app.Controllers.Corporation.ChatController.GetRooms)
		chat.PUT("/room/:roomID/block", app.Controllers.Corporation.ChatController.BlockRoom)
		chat.PUT("/room/:roomID/unblock", app.Controllers.Corporation.ChatController.UnBlockRoom)
	}

	blog := routerGroup.Group("/blog")
	{
		blog.POST("/create", app.Controllers.Corporation.BlogController.CreateDraftPost)
		blog.PUT("/:postID/edit", app.Controllers.Corporation.BlogController.EditPost)
		blog.PUT("/:postID/publish", app.Controllers.Corporation.BlogController.PublishPost)
		blog.PUT("/:postID/unpublish", app.Controllers.Corporation.BlogController.UnpublishPost)
		blog.DELETE("/", app.Controllers.Corporation.BlogController.DeletePost)
		blog.POST("/:postID/media", app.Controllers.Corporation.BlogController.AddPostMedia)
		blog.DELETE("/:postID/media/:mediaID", app.Controllers.Corporation.BlogController.DeletePostMedia)
		blog.GET("/list", app.Controllers.Corporation.BlogController.GetPosts)
		blog.GET("/:postID", app.Controllers.Corporation.BlogController.GetPost)
		blog.GET("/:postID/media/:mediaID", app.Controllers.Corporation.BlogController.GetPostMedia)
	}
}
