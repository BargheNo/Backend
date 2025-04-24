package maintenance

import (
	"time"

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

	pagination := controller.GetPagination(ctx, maintenanceController.pagination.DefaultPage, maintenanceController.pagination.DefaultPageSize)
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

func (maintenanceController *CorporationMaintenanceController) HandleMaintenanceRequest(ctx *gin.Context) {
	type handleMaintenanceRequestParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		RequestID     uint `json:"requestID" validate:"required"`
		Accept        bool `json:"accept" validate:"required"`
	}
	operatorID, _ := ctx.Get(maintenanceController.constants.Context.ID)
	params := controller.Validated[handleMaintenanceRequestParams](ctx)

	handleRequestInfo := maintenancedto.HandleRequest{
		CorporationID: params.CorporationID,
		OperatorID:    operatorID.(uint),
		RequestID:     params.RequestID,
		Accept:        params.Accept,
	}

	maintenanceController.maintenanceService.HandleRequest(handleRequestInfo)

	trans := controller.GetTranslator(ctx, maintenanceController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.maintenanceRequestHandled")
	controller.Response(ctx, 200, message, nil)
}

func (maintenanceController *CorporationMaintenanceController) AddMaintenanceRecord(ctx *gin.Context) {
	type addMaintenanceParams struct {
		CorporationID uint      `uri:"corporationID" validate:"required"`
		RequestID     uint      `json:"requestID" validate:"required"`
		Date          time.Time `json:"date" validate:"required"`
		Title         string    `json:"title" validate:"required"`
		Details       string    `json:"details" validate:"required"`
	}
	operatorID, _ := ctx.Get(maintenanceController.constants.Context.ID)
	params := controller.Validated[addMaintenanceParams](ctx)
	maintenanceRecordInfo := maintenancedto.AddMaintenanceRecordRequest{
		CorporationID: params.CorporationID,
		OperatorID:    operatorID.(uint),
		Date:          params.Date,
		Title:         params.Title,
		Details:       params.Details,
	}
	maintenanceController.maintenanceService.AddMaintenanceRecord(maintenanceRecordInfo)

	trans := controller.GetTranslator(ctx, maintenanceController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.addMaintenanceRecord")
	controller.Response(ctx, 200, message, nil)
}

func (maintenanceController *CorporationMaintenanceController) GetCorporationMaintenanceRecords(ctx *gin.Context) {
	type maintenanceRecordsParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	operatorID, _ := ctx.Get(maintenanceController.constants.Context.ID)
	params := controller.Validated[maintenanceRecordsParams](ctx)

	pagination := controller.GetPagination(ctx, maintenanceController.pagination.DefaultPage, maintenanceController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()
	listInfo := maintenancedto.CorporationMaintenanceListRequest{
		CorporationID: params.CorporationID,
		OperatorID:    operatorID.(uint),
		Offset:        offset,
		Limit:         limit,
	}

	requests := maintenanceController.maintenanceService.GetCorporationMaintenanceRecords(listInfo)
	controller.Response(ctx, 200, "success", requests)
}

func (maintenanceController *CorporationMaintenanceController) GetCorporationMaintenanceRecordsByPanel(ctx *gin.Context) {
	type maintenanceRecordsParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		PanelID       uint `uri:"panelID" validate:"required"`
	}
	operatorID, _ := ctx.Get(maintenanceController.constants.Context.ID)
	params := controller.Validated[maintenanceRecordsParams](ctx)

	pagination := controller.GetPagination(ctx, maintenanceController.pagination.DefaultPage, maintenanceController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()
	listInfo := maintenancedto.CorporationMaintenanceRecordByPanelRequest{
		CorporationID: params.CorporationID,
		OperatorID:    operatorID.(uint),
		PanelID:       params.PanelID,
		Offset:        offset,
		Limit:         limit,
	}

	requests := maintenanceController.maintenanceService.GetCorporationMaintenanceRecordsByPanel(listInfo)
	controller.Response(ctx, 200, "success", requests)
}
