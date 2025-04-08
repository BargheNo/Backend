package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	auth := routerGroup.Group("/auth")
	{
		auth.POST("/reset-password", app.Controllers.Customer.UserController.ResetPassword)
		// auth.POST("/verify/email", app.Controllers.General.UserController.VerifyEmail)
		// auth.POST("/register/complete", app.Controllers.General.UserController.CompleteRegister)
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
