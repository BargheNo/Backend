package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupMemberRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	corp := routerGroup.Group("/corp")
	{
		corp.GET("/installation_requests", app.Controllers.General.CorporationController.GetInstallationRequests)
	}
}
