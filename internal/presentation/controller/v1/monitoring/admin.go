package monitoring

import (
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/gin-gonic/gin"
)

type AdminMonitoringController struct {
	monitoringService service.MonitoringService
}

func NewAdminMonitoringController(monitoringService service.MonitoringService) *AdminMonitoringController {
	return &AdminMonitoringController{monitoringService: monitoringService}
}

func (c *AdminMonitoringController) Test(ctx *gin.Context) {}
