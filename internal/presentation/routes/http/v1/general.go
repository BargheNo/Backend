package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupGeneralRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	auth := routerGroup.Group("/auth")
	{
		auth.POST("/register/basic", app.Controllers.General.UserController.BasicRegister)
		auth.POST("/verify/phone", app.Controllers.General.UserController.VerifyPhone)
		auth.POST("/login", app.Controllers.General.UserController.Login)
		auth.POST("/forgot-password", app.Controllers.General.UserController.ForgotPassword)
		auth.POST("/confirm-otp", app.Controllers.General.UserController.ConfirmOTP)
		auth.POST("/refresh", app.Controllers.General.UserController.RefreshToken)
	}

	addresses := routerGroup.Group("/address")
	{
		addresses.GET("/province", app.Controllers.General.AddressController.GetProvince)
		addresses.GET("/province/:provinceID/city", app.Controllers.General.AddressController.GetProvinceCities)
	}

	contacts := routerGroup.Group("/contact")
	{
		contacts.GET("/types", app.Controllers.General.CorporationController.GetContactTypes)
	}

	notifications := routerGroup.Group("/notifications")
	{
		notifications.GET("/type", app.Controllers.General.NotificationController.GetContactTypes)
	}

	installations := routerGroup.Group("/installation")
	{
		requests := installations.Group("/request")
		{
			requests.GET("/status", app.Controllers.General.InstallationController.GetRequestStatuses)
			requests.GET("/building", app.Controllers.General.InstallationController.GetBuildingTypes)
		}
	}

	guarantees := routerGroup.Group("/guarantee")
	{
		guarantees.GET("/status", app.Controllers.Corporation.GuaranteeController.GetGuaranteeStatuses)
	}
}
