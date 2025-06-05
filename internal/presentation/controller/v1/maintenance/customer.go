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

// CHECKED
func (maintenanceController *CustomerMaintenanceController) GetMaintenanceUrgencyLevels(ctx *gin.Context) {
	levels := maintenanceController.maintenanceService.GetMaintenanceUrgencyLevels()

	controller.Response(ctx, 200, "", levels)
}

func (maintenanceController *CustomerMaintenanceController) GetMaintenanceStatuses(ctx *gin.Context) {
	statuses := maintenanceController.maintenanceService.GetMaintenanceRequestStatuses(enum.AgentTypeCustomer)

	controller.Response(ctx, 200, "", statuses)
}

// CHECKED
func (maintenanceController *CustomerMaintenanceController) CreateMaintenanceRequest(ctx *gin.Context) {
	type maintenanceRequestParams struct {
		PanelID          uint   `json:"panelID" validate:"required"`
		CorporationID    uint   `json:"corporationID" validate:"required"`
		Subject          string `json:"subject" validate:"required"`
		Description      string `json:"description" validate:"required"`
		UrgencyLevel     uint   `json:"urgencyLevel" validate:"required"`
		IsUsingGuarantee bool   `json:"isUsingGuarantee"`
	}
	params := controller.Validated[maintenanceRequestParams](ctx)
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)

	requestInfo := maintenancedto.CreateMaintenanceRequest{
		PanelID:          params.PanelID,
		OwnerID:          ownerID.(uint),
		CorporationID:    params.CorporationID,
		Subject:          params.Subject,
		Description:      params.Description,
		UrgencyLevel:     enum.UrgencyLevel(params.UrgencyLevel),
		IsUsingGuarantee: params.IsUsingGuarantee,
	}
	maintenanceController.maintenanceService.CreateMaintenanceRequest(requestInfo)

	trans := controller.GetTranslator(ctx, maintenanceController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.maintenanceRequest")
	controller.Response(ctx, 201, message, nil)
}

// CHECKED
func (maintenanceController *CustomerMaintenanceController) GetAllMaintenanceRequests(ctx *gin.Context) {
	type maintenanceRequestsParams struct {
		Status uint `form:"status" validate:"required"`
	}
	params := controller.Validated[maintenanceRequestsParams](ctx)
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)

	pagination := controller.GetPagination(ctx, maintenanceController.pagination.DefaultPage, maintenanceController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	listInfo := maintenancedto.CustomerMaintenanceListRequest{
		OwnerID: ownerID.(uint),
		Status:  params.Status,
		Offset:  offset,
		Limit:   limit,
	}
	requests := maintenanceController.maintenanceService.GetCustomerMaintenanceRequests(listInfo)

	controller.Response(ctx, 200, "", requests)
}

// CHECKED
func (maintenanceController *CustomerMaintenanceController) GetPanelMaintenanceRequests(ctx *gin.Context) {
	type maintenanceRequestsParams struct {
		PanelID uint `uri:"panelID" validate:"required"`
		Status  uint `form:"status" validate:"required"`
	}
	params := controller.Validated[maintenanceRequestsParams](ctx)
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)

	pagination := controller.GetPagination(ctx, maintenanceController.pagination.DefaultPage, maintenanceController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	listInfo := maintenancedto.CustomerPanelMaintenanceListRequest{
		OwnerID: ownerID.(uint),
		PanelID: params.PanelID,
		Status:  params.Status,
		Offset:  offset,
		Limit:   limit,
	}
	requests := maintenanceController.maintenanceService.GetCustomerPanelMaintenanceRequests(listInfo)

	controller.Response(ctx, 200, "", requests)
}

// CHECKED
func (maintenanceController *CustomerMaintenanceController) GetMaintenanceRequest(ctx *gin.Context) {
	type maintenanceRequestParams struct {
		RequestID uint `uri:"requestID" validate:"required"`
	}
	params := controller.Validated[maintenanceRequestParams](ctx)
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)

	maintenanceInfo := maintenancedto.CustomerMaintenanceRequest{
		OwnerID:   ownerID.(uint),
		RequestID: params.RequestID,
	}
	request := maintenanceController.maintenanceService.GetCustomerMaintenanceRequest(maintenanceInfo)

	controller.Response(ctx, 200, "", request)
}

// CHECKED
func (maintenanceController *CustomerMaintenanceController) UpdateMaintenanceRequest(ctx *gin.Context) {
	type updateMaintenanceRequestParams struct {
		RequestID        uint    `uri:"requestID" validate:"required"`
		Subject          *string `json:"subject"`
		Description      *string `json:"description"`
		UrgencyLevel     *uint   `json:"urgencyLevel"`
		IsUsingGuarantee *bool   `json:"isUsingGuarantee"`
	}
	params := controller.Validated[updateMaintenanceRequestParams](ctx)
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)

	requestInfo := maintenancedto.UpdateCustomerRequest{
		OwnerID:          ownerID.(uint),
		RequestID:        params.RequestID,
		Subject:          params.Subject,
		Description:      params.Description,
		UrgencyLevel:     params.UrgencyLevel,
		IsUsingGuarantee: params.IsUsingGuarantee,
	}
	maintenanceController.maintenanceService.UpdateMaintenanceRequest(requestInfo)

	trans := controller.GetTranslator(ctx, maintenanceController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateMaintenanceRequest")
	controller.Response(ctx, 201, message, nil)
}

// CHECKED
func (maintenanceController *CustomerMaintenanceController) CancelMaintenanceRequest(ctx *gin.Context) {
	type cancelMaintenanceRequestParams struct {
		RequestID uint `uri:"requestID" validate:"required"`
	}
	params := controller.Validated[cancelMaintenanceRequestParams](ctx)
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)

	maintenanceInfo := maintenancedto.CustomerMaintenanceRequest{
		OwnerID:   ownerID.(uint),
		RequestID: params.RequestID,
	}
	maintenanceController.maintenanceService.CancelMaintenanceRequest(maintenanceInfo)

	trans := controller.GetTranslator(ctx, maintenanceController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.cancelMaintenanceRequest")
	controller.Response(ctx, 201, message, nil)
}

// CHECKED
func (maintenanceController *CustomerMaintenanceController) ApproveMaintenanceRecord(ctx *gin.Context) {
	type cancelMaintenanceRequestParams struct {
		RequestID uint `uri:"requestID" validate:"required"`
	}
	params := controller.Validated[cancelMaintenanceRequestParams](ctx)
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)

	maintenanceInfo := maintenancedto.CustomerMaintenanceRequest{
		OwnerID:   ownerID.(uint),
		RequestID: params.RequestID,
	}
	maintenanceController.maintenanceService.ApproveMaintenanceRecord(maintenanceInfo)

	trans := controller.GetTranslator(ctx, maintenanceController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.approveMaintenanceRecord")
	controller.Response(ctx, 201, message, nil)
}
