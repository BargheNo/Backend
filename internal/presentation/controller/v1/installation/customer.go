package installation

import (
	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerInstallationController struct {
	constants           *bootstrap.Constants
	pagination          *bootstrap.Pagination
	installationService usecase.InstallationService
}

func NewCustomerInstallationController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	installationService usecase.InstallationService,
) *CustomerInstallationController {
	return &CustomerInstallationController{
		constants:           constants,
		pagination:          pagination,
		installationService: installationService,
	}
}

func (installationController *CustomerInstallationController) CreateInstallationRequest(ctx *gin.Context) {
	type installationRequestParams struct {
		Name          string  `json:"name" validate:"required"`
		Area          uint    `json:"area"`
		Power         uint    `json:"power" validate:"required"`
		MaxCost       float64 `json:"maxCost"`
		BuildingType  uint    `json:"buildingType" validate:"required"`
		Description   string  `json:"description"`
		ProvinceID    uint    `json:"provinceID" validate:"required"`
		CityID        uint    `json:"cityID" validate:"required"`
		StreetAddress string  `json:"streetAddress" validate:"required"`
		PostalCode    string  `json:"postalCode" validate:"required"`
		HouseNumber   string  `json:"houseNumber" validate:"required"`
		Unit          uint    `json:"unit" validate:"required"`
	}
	params := controller.Validated[installationRequestParams](ctx)
	ownerID, _ := ctx.Get(installationController.constants.Context.ID)
	requestInfo := installationdto.NewInstallationRequest{
		OwnerID:      ownerID.(uint),
		Name:         params.Name,
		Area:         params.Area,
		Power:        params.Power,
		MaxCost:      params.MaxCost,
		BuildingType: params.BuildingType,
		Description:  params.Description,
		Address: addressdto.CreateAddressRequest{
			ProvinceID:    params.ProvinceID,
			CityID:        params.CityID,
			StreetAddress: params.StreetAddress,
			PostalCode:    params.PostalCode,
			HouseNumber:   params.HouseNumber,
			Unit:          params.Unit,
		},
	}
	if err := installationController.installationService.CreateInstallationRequest(requestInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.installationRequest")
	controller.Response(ctx, 201, message, nil)
}

func (installationController *CustomerInstallationController) GetInstallationRequests(ctx *gin.Context) {
	type getRequestsParams struct {
		Status   uint `form:"status"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
		SortBy   uint `form:"sortBy"`
		Asc      bool `form:"asc"`
	}
	params := controller.Validated[getRequestsParams](ctx)

	ownerID, _ := ctx.Get(installationController.constants.Context.ID)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, installationController.pagination.DefaultPage, installationController.pagination.DefaultPageSize)

	listInfo := installationdto.CustomerRequestsListRequest{
		OwnerID: ownerID.(uint),
		Status:  params.Status,
		Offset:  offset,
		Limit:   limit,
		SortBy:  params.SortBy,
		Asc:     params.Asc,
	}
	requests, count, err := installationController.installationService.GetOwnerInstallationRequests(listInfo)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(requests, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (installationController *CustomerInstallationController) GetInstallationRequest(ctx *gin.Context) {
	type installationRequestParams struct {
		RequestID uint `uri:"requestID" validate:"required"`
	}
	params := controller.Validated[installationRequestParams](ctx)
	ownerID, _ := ctx.Get(installationController.constants.Context.ID)

	requestInfo := installationdto.GetOwnerRequest{
		InstallationID: params.RequestID,
		OwnerID:        ownerID.(uint),
	}
	installationRequest, err := installationController.installationService.GetOwnerInstallationRequest(requestInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", installationRequest)
}

func (installationController *CustomerInstallationController) CancelInstallationRequest(ctx *gin.Context) {
	type installationRequestParams struct {
		RequestID uint `uri:"requestID" validate:"required"`
	}
	params := controller.Validated[installationRequestParams](ctx)
	ownerID, _ := ctx.Get(installationController.constants.Context.ID)

	requestInfo := installationdto.ChangeRequestStatusRequest{
		RequestID: params.RequestID,
		Status:    enum.InstallationRequestStatusCancelled,
		OwnerID:   ownerID.(uint),
	}
	if err := installationController.installationService.ChangeInstallationRequestStatus(requestInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.cancelInstallationRequest")
	controller.Response(ctx, 201, message, nil)
}

func (installationController *CustomerInstallationController) GetCustomerPanels(ctx *gin.Context) {
	type getPanelsParams struct {
		Status   uint `form:"status"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
		SortBy   uint `form:"sortBy"`
		Asc      bool `form:"asc"`
	}
	params := controller.Validated[getPanelsParams](ctx)

	ownerId, _ := ctx.Get(installationController.constants.Context.ID)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, installationController.pagination.DefaultPage, installationController.pagination.DefaultPageSize)

	listInfo := installationdto.CustomerPanelListRequest{
		OwnerID: ownerId.(uint),
		Status:  params.Status,
		Offset:  offset,
		Limit:   limit,
		SortBy:  params.SortBy,
		Asc:     params.Asc,
	}
	panels, count, err := installationController.installationService.GetCustomerPanels(listInfo)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(panels, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (installationController *CustomerInstallationController) GetCustomerPanel(ctx *gin.Context) {
	type getPanelParams struct {
		PanelID uint `uri:"panelID" validate:"required"`
	}
	params := controller.Validated[getPanelParams](ctx)
	ownerID, _ := ctx.Get(installationController.constants.Context.ID)

	panelInfo := installationdto.GetOwnerRequest{
		InstallationID: params.PanelID,
		OwnerID:        ownerID.(uint),
	}
	panels, err := installationController.installationService.GetCustomerPanel(panelInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", panels)
}

func (installationController *CustomerInstallationController) GetPanelGuaranteeViolation(ctx *gin.Context) {
	type getPanelParams struct {
		PanelID uint `uri:"panelID" validate:"required"`
	}
	params := controller.Validated[getPanelParams](ctx)
	ownerID, _ := ctx.Get(installationController.constants.Context.ID)

	violationInfo := installationdto.GetCustomerGuaranteeViolationRequest{
		OwnerID: ownerID.(uint),
		PanelID: params.PanelID,
	}
	panels, err := installationController.installationService.GetCustomerPanelGuaranteeViolation(violationInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", panels)
}
