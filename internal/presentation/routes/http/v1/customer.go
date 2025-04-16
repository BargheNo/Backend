package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	profile := routerGroup.Group("/profile")
	{
		profile.PUT("/password", app.Controllers.Customer.UserController.ResetPassword)
		profile.POST("/complete", app.Controllers.Customer.UserController.CompleteRegister)
		profile.POST("/verify/email", app.Controllers.General.UserController.VerifyEmail)
		// auth.POST("/verify/email", app.Controllers.General.UserController.VerifyEmail)
	}

	corp := routerGroup.Group("/corp")
	{
		corp.POST("/register", app.Controllers.Customer.CorporationController.Register)
		corp.GET("/list", app.Controllers.Customer.CorporationController.GetCorporations)
	}

	orders := routerGroup.Group("/installation")
	{
		requests := orders.Group("/request")
		{
			requests.POST("", app.Controllers.Customer.InstallationController.CreateInstallationRequest)
			requests.GET("", app.Controllers.Customer.InstallationController.GetOwnerInstallationRequests)
			requests.GET("/:requestID", app.Controllers.Customer.InstallationController.GetInstallationRequest)
			requests.GET("/:requestID/bids", app.Controllers.Customer.BidController.GetBids)
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
	}

	notification := routerGroup.Group("/notification")
	{
		notification.POST("/:notificationID/read", app.Controllers.Customer.NotificationController.MarkAsRead)
		notification.GET("", app.Controllers.Customer.NotificationController.GetUserNotifications)
		notification.GET("/setting", app.Controllers.Customer.NotificationController.GetUserNotificationSettings)
		notification.PUT("/setting/:settingID", app.Controllers.Customer.NotificationController.UpdateSettings)
	}

	panels := routerGroup.Group("/panels")
	{
		panels.GET("/list", app.Controllers.Customer.InstallationController.GetCustomerPanels)
	}

	maintenance := routerGroup.Group("/maintenance")
	{
		requests := maintenance.Group("/request")
		{
			requests.POST("/", app.Controllers.Customer.MaintenanceController.CreateMaintenanceRequest)
			requests.GET("/list", app.Controllers.Customer.MaintenanceController.GetCustomerMaintenanceRequests)
		}

		records := maintenance.Group("/record")
		{
			records.GET("/list", app.Controllers.Customer.MaintenanceController.GetMaintenanceRecords)
			records.GET("/list/:panelID", app.Controllers.Customer.MaintenanceController.GetCustomerMaintenanceRequestsByPanelID)
		}
	}
}
