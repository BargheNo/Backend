package monitoring

import (
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/gin-gonic/gin"
)

type AdminMonitoringController struct {
	monitoringService usecase.MonitoringService
}

func NewAdminMonitoringController(monitoringService usecase.MonitoringService) *AdminMonitoringController {
	return &AdminMonitoringController{monitoringService: monitoringService}
}

func (c *AdminMonitoringController) Test(ctx *gin.Context) {}
