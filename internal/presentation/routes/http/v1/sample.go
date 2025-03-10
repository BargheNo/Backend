package httpv1

import (
	"github.com/BargheNo/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupSampleRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	routerGroup.POST("/sample", app.Controllers.SampleController.SampleCreate)
}
