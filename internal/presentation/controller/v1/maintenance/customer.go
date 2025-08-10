package maintenance

import (
	"github.com/BargheNo/Backend/bootstrap"
	maintenancedto "github.com/BargheNo/Backend/internal/application/dto/maintenance"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerMaintenanceController struct {
	constants          *bootstrap.Constants
	pagination         *bootstrap.Pagination
	maintenanceService usecase.MaintenanceService
}

func NewCustomerMaintenanceController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	maintenanceService usecase.MaintenanceService,
) *CustomerMaintenanceController {
	return &CustomerMaintenanceController{
		constants:          constants,
		pagination:         pagination,
		maintenanceService: maintenanceService,
	}
}

func (maintenanceController *CustomerMaintenanceController) GetMaintenanceUrgencyLevels(ctx *gin.Context) {
	levels := maintenanceController.maintenanceService.GetMaintenanceUrgencyLevels()

	controller.Response(ctx, 200, "", levels)
}

func (maintenanceController *CustomerMaintenanceController) GetMaintenanceStatuses(ctx *gin.Context) {
	statuses := maintenanceController.maintenanceService.GetMaintenanceRequestStatuses(enum.AgentTypeCustomer)

	controller.Response(ctx, 200, "", statuses)
}

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
	if err := maintenanceController.maintenanceService.CreateMaintenanceRequest(requestInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, maintenanceController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.maintenanceRequest")
	controller.Response(ctx, 201, message, nil)
}

func (maintenanceController *CustomerMaintenanceController) GetAllMaintenanceRequests(ctx *gin.Context) {
	type maintenanceRequestsParams struct {
		Status   uint `form:"status"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
		SortBy   uint `form:"sortBy"`
		Asc      bool `form:"asc"`
	}
	params := controller.Validated[maintenanceRequestsParams](ctx)
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, maintenanceController.pagination.DefaultPage, maintenanceController.pagination.DefaultPageSize)

	listInfo := maintenancedto.CustomerMaintenanceListRequest{
		OwnerID: ownerID.(uint),
		Status:  params.Status,
		Offset:  offset,
		Limit:   limit,
		SortBy:  params.SortBy,
		Asc:     params.Asc,
	}
	requests, count, err := maintenanceController.maintenanceService.GetCustomerMaintenanceRequests(listInfo)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(requests, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (maintenanceController *CustomerMaintenanceController) GetPanelMaintenanceRequests(ctx *gin.Context) {
	type maintenanceRequestsParams struct {
		PanelID  uint `uri:"panelID" validate:"required"`
		Status   uint `form:"status"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
		SortBy   uint `form:"sortBy"`
		Asc      bool `form:"asc"`
	}
	params := controller.Validated[maintenanceRequestsParams](ctx)
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, maintenanceController.pagination.DefaultPage, maintenanceController.pagination.DefaultPageSize)

	listInfo := maintenancedto.CustomerPanelMaintenanceListRequest{
		OwnerID: ownerID.(uint),
		PanelID: params.PanelID,
		Status:  params.Status,
		Offset:  offset,
		Limit:   limit,
		SortBy:  params.SortBy,
		Asc:     params.Asc,
	}
	requests, count, err := maintenanceController.maintenanceService.GetCustomerPanelMaintenanceRequests(listInfo)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(requests, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

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
	request, err := maintenanceController.maintenanceService.GetCustomerMaintenanceRequest(maintenanceInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", request)
}

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
	if err := maintenanceController.maintenanceService.UpdateMaintenanceRequest(requestInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, maintenanceController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateMaintenanceRequest")
	controller.Response(ctx, 201, message, nil)
}

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
	if err := maintenanceController.maintenanceService.CancelMaintenanceRequest(maintenanceInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, maintenanceController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.cancelMaintenanceRequest")
	controller.Response(ctx, 201, message, nil)
}

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
	if err := maintenanceController.maintenanceService.ApproveMaintenanceRecord(maintenanceInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, maintenanceController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.approveMaintenanceRecord")
	controller.Response(ctx, 201, message, nil)
}
