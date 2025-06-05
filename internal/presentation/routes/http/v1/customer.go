package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	profile := routerGroup.Group("/profile")
	{
		profile.GET("", app.Controllers.Customer.UserController.GetMyProfile)
		profile.PUT("/password", app.Controllers.Customer.UserController.ResetPassword)
		profile.POST("/complete", app.Controllers.Customer.UserController.CompleteRegister)
		profile.POST("/verify/email", app.Controllers.Customer.UserController.VerifyEmail)
		profile.PUT("", app.Controllers.Customer.UserController.UpdateProfile)
	}

	corps := routerGroup.Group("/corps")
	{
		registration := corps.Group("/registration")
		{
			registration.POST("/basic", app.Controllers.Customer.CorporationController.Register)
			corpsSubgroup := registration.Group("/:corporationID")
			{
				corpsSubgroup.PUT("/basic", app.Controllers.Customer.CorporationController.UpdateRegister)
				corpsSubgroup.POST("/contacts", app.Controllers.Customer.CorporationController.AddContactInformation)
				corpsSubgroup.DELETE("/contacts/:contactID", app.Controllers.Customer.CorporationController.DeleteContactInformation)
				corpsSubgroup.POST("/address", app.Controllers.Customer.CorporationController.AddAddress)
				corpsSubgroup.DELETE("/address/:addressID", app.Controllers.Customer.CorporationController.DeleteAddress)
				corpsSubgroup.PUT("/certificates", app.Controllers.Customer.CorporationController.SubmitCertificateFiles)
				corpsSubgroup.GET("", app.Controllers.Customer.CorporationController.GetCorporationPrivateDetails)
			}
		}
	}

	installations := routerGroup.Group("/installation")
	{
		requests := installations.Group("/request")
		{
			requests.POST("", app.Controllers.Customer.InstallationController.CreateInstallationRequest)
			requests.GET("", app.Controllers.Customer.InstallationController.GetInstallationRequests)
			requestSubGroup := requests.Group("/:requestID")
			{
				requestSubGroup.GET("", app.Controllers.Customer.InstallationController.GetInstallationRequest)
				requestSubGroup.PUT("/cancel", app.Controllers.Customer.InstallationController.CancelInstallationRequest)

				bids := requestSubGroup.Group("/bid")
				{
					bids.GET("", app.Controllers.Customer.BidController.GetBids)
					bidsSubGroup := bids.Group("/:bidID")
					{
						bidsSubGroup.GET("", app.Controllers.Customer.BidController.GetBid)
						bidsSubGroup.POST("/accept", app.Controllers.Customer.BidController.AcceptBid)
						bidsSubGroup.POST("/reject", app.Controllers.Customer.BidController.RejectBid)
					}
				}
			}
		}

		panels := installations.Group("/panel")
		{
			panels.GET("", app.Controllers.Customer.InstallationController.GetCustomerPanels)
			panelsSubGroup := panels.Group("/:panelID")
			{
				panelsSubGroup.GET("", app.Controllers.Customer.InstallationController.GetCustomerPanel)
				panelsSubGroup.GET("/guarantee/violation", app.Controllers.Customer.InstallationController.GetPanelGuaranteeViolation)
				panelsSubGroup.GET("/maintenance", app.Controllers.Customer.MaintenanceController.GetPanelMaintenanceRequests)
			}
		}
	}

	maintenances := routerGroup.Group("/maintenance/request")
	{
		maintenances.GET("/level", app.Controllers.Customer.MaintenanceController.GetMaintenanceUrgencyLevels)
		maintenances.POST("", app.Controllers.Customer.MaintenanceController.CreateMaintenanceRequest)
		maintenances.GET("", app.Controllers.Customer.MaintenanceController.GetAllMaintenanceRequests)
		requestsSubGroup := maintenances.Group("/:requestID")
		{
			requestsSubGroup.GET("", app.Controllers.Customer.MaintenanceController.GetMaintenanceRequest)
			requestsSubGroup.PUT("", app.Controllers.Customer.MaintenanceController.UpdateMaintenanceRequest)
			requestsSubGroup.PUT("/cancel", app.Controllers.Customer.MaintenanceController.CancelMaintenanceRequest)
			requestsSubGroup.PUT("record/approve", app.Controllers.Customer.MaintenanceController.ApproveMaintenanceRecord)
			// records := requestsSubGroup.Group("/record")
			// {
			// 	// records.GET("") // i think we can ignore this for now :)
			// 	// records.PUT("/:recordID/approve") // can ignore recordID part :)
			// 	// records.PUT("/approve")
			// }
		}
	}

	addresses := routerGroup.Group("/address")
	{
		addresses.POST("", app.Controllers.Customer.AddressController.CreateUserAddress)
		addresses.GET("", app.Controllers.Customer.AddressController.GetCustomerAddresses)
	}

	chat := routerGroup.Group("/chat")
	{
		chat.POST("/room/:corporationID", app.Controllers.Customer.ChatController.CreateOrGetRoom)
		chat.GET("/room", app.Controllers.Customer.ChatController.GetUserRooms)
		chat.GET("/room/:roomID/messages", app.Controllers.Customer.ChatController.GetMessages)
		chat.PUT("/room/:roomID/block", app.Controllers.Customer.ChatController.BlockRoom)
		chat.PUT("/room/:roomID/unblock", app.Controllers.Customer.ChatController.UnBlockRoom)
	}

	notification := routerGroup.Group("/notifications")
	{
		notification.POST("/:notificationID/read", app.Controllers.Customer.NotificationController.MarkAsRead)
		notification.GET("", app.Controllers.Customer.NotificationController.GetUserNotifications)
		notification.GET("/setting", app.Controllers.Customer.NotificationController.GetUserNotificationSettings)
		notification.PUT("/setting/:settingID", app.Controllers.Customer.NotificationController.UpdateSettings)
	}

	tickets := routerGroup.Group("/ticket")
	{
		tickets.POST("", app.Controllers.Customer.TicketController.CreateTicket)
		tickets.GET("/list", app.Controllers.Customer.TicketController.GetTickets)
		ticketSubGroup := tickets.Group("/:ticketID/comments")
		{
			ticketSubGroup.GET("", app.Controllers.Customer.TicketController.GetComments)
			ticketSubGroup.POST("", app.Controllers.Customer.TicketController.CreateComment)
		}
	}

	reports := routerGroup.Group("/report")
	{
		maintenanceReports := reports.Group("/maintenance")
		{
			maintenanceReports.POST("/:recordID", app.Controllers.Customer.ReportController.CreateMaintenanceReport)
		}
		panelReports := reports.Group("/panel")
		{
			panelReports.POST("/:panelID", app.Controllers.Customer.ReportController.CreatePanelReport)
		}
	}

	blog := routerGroup.Group("/blog")
	{
		blog.POST("/:postID/like", app.Controllers.Customer.BlogController.LikePost)
		blog.DELETE("/:postID/like", app.Controllers.Customer.BlogController.UnlikePost)
	}
}
