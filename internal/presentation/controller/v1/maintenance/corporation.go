package maintenance

import (
	"strconv"

	bootstrap "github.com/BargheNo/Backend/bootstrap"

	maintenancedto "github.com/BargheNo/Backend/internal/application/dto/maintenance"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CorporationMaintenanceController struct {
	constants          *bootstrap.Constants
	pagination         *bootstrap.Pagination
	maintenanceService service.MaintenanceService
}

func NewCorporationMaintenanceController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	maintenanceService service.MaintenanceService,
) *CorporationMaintenanceController {
	return &CorporationMaintenanceController{
		constants:          constants,
		pagination:         pagination,
		maintenanceService: maintenanceService,
	}
}

func (maintenanceController *CorporationMaintenanceController) GetMaintenanceRequests(ctx *gin.Context) {
	type maintenanceRequestParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	operatorID, _ := ctx.Get(maintenanceController.constants.Context.ID)
	params := controller.Validated[maintenanceRequestParams](ctx)

	defaultPage, err := strconv.Atoi(maintenanceController.pagination.DefaultPage)
	if err != nil {
		defaultPage = 1
	}
	defaultPageSize, err := strconv.Atoi(maintenanceController.pagination.DefaultPageSize)
	if err != nil {
		defaultPageSize = 10
	}
	pagination := controller.GetPagination(ctx, defaultPage, defaultPageSize)
	offset, limit := pagination.GetOffsetLimit()
	listInfo := maintenancedto.CorporationMaintenanceListRequest{
		CorporationID: params.CorporationID,
		OperatorID:    operatorID.(uint),
		Offset:        offset,
		Limit:         limit,
	}

	requests := maintenanceController.maintenanceService.GetCorporationMaintenanceRequests(listInfo)
	controller.Response(ctx, 200, "success", requests)
}
