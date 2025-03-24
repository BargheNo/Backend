package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	auth := routerGroup.Group("/auth")
	{
		corp := auth.Group("/corporation")
		{
			corp.POST("/reset-password", app.Controllers.Customer.CorporationController.ChangePassword)
		}
		auth.POST("/reset-password", app.Controllers.Customer.UserController.ResetPassword)

		// auth.POST("/verify/email", app.Controllers.General.UserController.VerifyEmail)
		// auth.POST("/register/complete", app.Controllers.General.UserController.CompleteRegister)
	}

	orders := routerGroup.Group("/installation")
	{
		orders.POST("/request", app.Controllers.Customer.InstallationController.CreateInstallationRequest)
		orders.GET("/request", app.Controllers.Customer.InstallationController.GetOwnerInstallationRequests)
		orders.GET("/request/:requestID", app.Controllers.Customer.InstallationController.GetInstallationRequest)
	}

	addresses := routerGroup.Group("/address")
	{
		addresses.POST("/user", app.Controllers.Customer.AddressController.CreateUserAddress)
		addresses.GET("/user", app.Controllers.Customer.AddressController.GetCustomerAddresses)
		addresses.POST("/corp", app.Controllers.Customer.AddressController.CreateCorporationAddress)
		addresses.DELETE("/corp", app.Controllers.Customer.AddressController.DeleteCorporationAddress)
	}

	bids := routerGroup.Group("/bids")
	{
		bids.POST("/set", app.Controllers.Customer.CorporationController.SetBid)
		bids.POST("/cancel", app.Controllers.Customer.CorporationController.CancelBid)
		bids.GET("/list", app.Controllers.Customer.CorporationController.GetBids)
	}

	corp := routerGroup.Group("/corp")
	{
		corp.GET("/info", app.Controllers.Customer.CorporationController.GetCorporationInfo)
		corp.POST("/contact-info", app.Controllers.Customer.CorporationController.UpdateContactInfo)
	}
}
