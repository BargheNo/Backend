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

	orders := routerGroup.Group("/installation")
	{
		requests := orders.Group("/request")
		{
			requests.POST("", app.Controllers.Customer.InstallationController.CreateInstallationRequest)
			requests.GET("", app.Controllers.Customer.InstallationController.GetOwnerInstallationRequests)
			requests.GET("/:requestID", app.Controllers.Customer.InstallationController.GetInstallationRequest)
		}
	}

	addresses := routerGroup.Group("/address")
	{
		addresses.POST("/user", app.Controllers.Customer.AddressController.CreateUserAddress)
		addresses.GET("/user", app.Controllers.Customer.AddressController.GetCustomerAddresses)
	}
}
