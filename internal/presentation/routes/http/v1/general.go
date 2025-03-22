package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupGeneralRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	auth := routerGroup.Group("/auth")
	{
		auth.POST("/register/basic", app.Controllers.General.UserController.BasicRegister)
		// auth.POST("/register/complete", app.Controllers.General.UserController.CompleteRegister)
		auth.POST("/verify/phone", app.Controllers.General.UserController.VerifyPhone)
		auth.POST("/verify/email", app.Controllers.General.UserController.VerifyEmail)
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
}
