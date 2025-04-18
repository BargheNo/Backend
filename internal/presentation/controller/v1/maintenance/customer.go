package maintenance

import (
	"github.com/BargheNo/Backend/bootstrap"
	maintenancedto "github.com/BargheNo/Backend/internal/application/dto/maintenance"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerMaintenanceController struct {
	constants          *bootstrap.Constants
	pagination         *bootstrap.Pagination
	maintenanceService service.MaintenanceService
}

func NewCustomerMaintenanceController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	maintenanceService service.MaintenanceService,
) *CustomerMaintenanceController {
	return &CustomerMaintenanceController{
		constants:          constants,
		pagination:         pagination,
		maintenanceService: maintenanceService,
	}
}

func (maintenanceController *CustomerMaintenanceController) CreateMaintenanceRequest(ctx *gin.Context) {
	type maintenanceRequestParams struct {
		PanelID       uint   `json:"panelID" validate:"required"`
		CorporationID uint   `uri:"corporationID" validate:"required"`
		Subject       string `json:"subject" validate:"required"`
		Description   string `json:"description" validate:"required"`
		UrgencyLevel  uint   `json:"urgencyLevel" validate:"required"`
	}
	params := controller.Validated[maintenanceRequestParams](ctx)
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)
	requestInfo := maintenancedto.NewMaintenanceRequest{
		PanelID:       params.PanelID,
		OwnerID:       ownerID.(uint),
		CorporationID: params.CorporationID,
		Subject:       params.Subject,
		Description:   params.Description,
		UrgencyLevel:  enum.UrgencyLevel(params.UrgencyLevel),
	}

	maintenanceController.maintenanceService.CreateMaintenanceRequest(requestInfo)

	trans := controller.GetTranslator(ctx, maintenanceController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.maintenanceRequest")
	controller.Response(ctx, 201, message, nil)
}

func (maintenanceController *CustomerMaintenanceController) GetCustomerMaintenanceRequests(ctx *gin.Context) {
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)
	params := controller.GetPagination(ctx, maintenanceController.pagination.DefaultPage, maintenanceController.pagination.DefaultPageSize)
	offset, limit := params.GetOffsetLimit()
	listInfo := maintenancedto.MaintenanceListRequest{
		OwnerID: ownerID.(uint),
		Offset:  offset,
		Limit:   limit,
	}
	requests := maintenanceController.maintenanceService.GetCustomerMaintenanceRequests(listInfo)
	controller.Response(ctx, 200, "", requests)
}

func (maintenanceController *CustomerMaintenanceController) GetMaintenanceRecords(ctx *gin.Context) {
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)
	params := controller.GetPagination(ctx, maintenanceController.pagination.DefaultPage, maintenanceController.pagination.DefaultPageSize)
	offset, limit := params.GetOffsetLimit()
	listInfo := maintenancedto.MaintenanceListRequest{
		OwnerID: ownerID.(uint),
		Offset:  offset,
		Limit:   limit,
	}
	records := maintenanceController.maintenanceService.GetCustomerMaintenanceRecords(listInfo)
	controller.Response(ctx, 200, "", records)
}

func (maintenanceController *CustomerMaintenanceController) GetCustomerMaintenanceRequestsByPanelID(ctx *gin.Context) {
	type maintenanceRecordsParams struct {
		PanelID uint `uri:"panelID" validate:"required"`
	}
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)
	params := controller.Validated[maintenanceRecordsParams](ctx)
	pagination := controller.GetPagination(ctx, maintenanceController.pagination.DefaultPage, maintenanceController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()
	listInfo := maintenancedto.CustomerMaintenanceRecordByPanelRequest{
		OwnerID: ownerID.(uint),
		PanelID: params.PanelID,
		Offset:  offset,
		Limit:   limit,
	}

	requests := maintenanceController.maintenanceService.GetCustomerMaintenanceRecordsByPanel(listInfo)
	controller.Response(ctx, 200, "", requests)
}
