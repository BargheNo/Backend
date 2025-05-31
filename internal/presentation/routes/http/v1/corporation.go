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

	blog := routerGroup.Group(":corporationID/blog")
	{
		blog.POST("/create", app.Controllers.Corporation.BlogController.CreateDraftPost)
		blog.PUT("/:postID/edit", app.Controllers.Corporation.BlogController.EditPost)
		blog.PUT("/:postID/publish", app.Controllers.Corporation.BlogController.PublishPost)
		blog.PUT("/:postID/unpublish", app.Controllers.Corporation.BlogController.UnpublishPost)
		blog.DELETE("/", app.Controllers.Corporation.BlogController.DeletePost)
		blog.POST("/:postID/media", app.Controllers.Corporation.BlogController.AddPostMedia)
		blog.DELETE("/:postID/media/:mediaID", app.Controllers.Corporation.BlogController.DeletePostMedia)
		blog.GET("/list", app.Controllers.Corporation.BlogController.GetPosts)
		// blog.GET("/:postID/media/:mediaID", app.Controllers.Corporation.BlogController.GetPostMedia)
	}
}
