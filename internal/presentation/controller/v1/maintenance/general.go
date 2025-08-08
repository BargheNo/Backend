package maintenance

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralMaintenanceController struct {
	constants          *bootstrap.Constants
	maintenanceService usecase.MaintenanceService
}

func NewGeneralMaintenanceController(
	constants *bootstrap.Constants,
	maintenanceService usecase.MaintenanceService,

) *GeneralMaintenanceController {
	return &GeneralMaintenanceController{
		constants:          constants,
		maintenanceService: maintenanceService,
	}
}

func (maintenanceController *GeneralMaintenanceController) GetSortableFields(ctx *gin.Context) {
	columns := maintenanceController.maintenanceService.GetMaintenanceSortableColumns()
	controller.Response(ctx, 200, "", columns)
}
