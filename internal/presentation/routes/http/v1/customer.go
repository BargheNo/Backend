package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	auth := routerGroup.Group("/auth")
	{
		auth.POST("/reset-password", app.Controllers.General.UserController.ResetPassword)
		auth.POST("/verify/email", app.Controllers.General.UserController.VerifyEmail)
		// auth.POST("/register/complete", app.Controllers.General.UserController.CompleteRegister)
	}
}
