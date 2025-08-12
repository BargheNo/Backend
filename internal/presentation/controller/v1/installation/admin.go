package installation

import (
	"github.com/BargheNo/Backend/bootstrap"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminInstallationController struct {
	constants           *bootstrap.Constants
	pagination          *bootstrap.Pagination
	installationService usecase.InstallationService
}

func NewAdminInstallationController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	installationService usecase.InstallationService,
) *AdminInstallationController {
	return &AdminInstallationController{
		constants:           constants,
		pagination:          pagination,
		installationService: installationService,
	}
}

func (installationController *AdminInstallationController) GetInstallationRequests(ctx *gin.Context) {
	type getRequestsParams struct {
		Status   uint `form:"status"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
		SortBy   uint `form:"sortBy"`
		Asc      bool `form:"asc"`
	}
	params := controller.Validated[getRequestsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, installationController.pagination.DefaultPage, installationController.pagination.DefaultPageSize)

	listInfo := installationdto.AdminInstallationListRequest{
		Status: params.Status,
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}
	requests, count, err := installationController.installationService.GetInstallationRequestsByAdmin(listInfo)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(requests, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (installationController *AdminInstallationController) SearchInstallationRequests(ctx *gin.Context) {
	type searchInstallationRequestsParams struct {
		Query    string `form:"query" validate:"required"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		SortBy   uint   `form:"sortBy"`
		Asc      bool   `form:"asc"`
	}
	params := controller.Validated[searchInstallationRequestsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, installationController.pagination.DefaultPage, installationController.pagination.DefaultPageSize)

	request := installationdto.SearchInstallationRequestsRequest{
		Query:  params.Query,
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}

	requests, count, err := installationController.installationService.SearchInstallationRequests(request)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(requests, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (installationController *AdminInstallationController) GetInstallationRequest(ctx *gin.Context) {
	type installationRequestParams struct {
		RequestID uint `uri:"requestID" validate:"required"`
	}
	params := controller.Validated[installationRequestParams](ctx)

	installationRequest, err := installationController.installationService.GetPublicInstallationRequest(params.RequestID)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", installationRequest)
}

func (installationController *AdminInstallationController) UpdateInstallationRequest(ctx *gin.Context) {
	type installationRequestParams struct {
		RequestID    uint     `uri:"requestID" validate:"required"`
		Name         *string  `json:"name"`
		Area         *uint    `json:"area"`
		Power        *uint    `json:"power"`
		MaxCost      *float64 `json:"maxCost"`
		BuildingType *uint    `json:"buildingType"`
		Status       *uint    `json:"status"`
		Description  *string  `json:"description"`
	}
	params := controller.Validated[installationRequestParams](ctx)

	requestInfo := installationdto.UpdateInstallationRequest{
		RequestID:    params.RequestID,
		Name:         params.Name,
		Area:         params.Area,
		Power:        params.Power,
		MaxCost:      params.MaxCost,
		BuildingType: params.BuildingType,
		Status:       params.Status,
		Description:  params.Description,
	}
	if err := installationController.installationService.UpdateInstallationRequestByAdmin(requestInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateInstallationRequest")
	controller.Response(ctx, 201, message, nil)
}

func (installationController *AdminInstallationController) DeleteInstallationRequest(ctx *gin.Context) {
	type installationRequestParams struct {
		RequestID uint `uri:"requestID" validate:"required"`
	}
	params := controller.Validated[installationRequestParams](ctx)

	if err := installationController.installationService.DeleteInstallationRequest(params.RequestID); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteInstallationRequest")
	controller.Response(ctx, 200, message, nil)
}

func (installationController *AdminInstallationController) GetPanels(ctx *gin.Context) {
	type getPanelsParams struct {
		Status   uint `form:"status"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
		SortBy   uint `form:"sortBy"`
		Asc      bool `form:"asc"`
	}
	params := controller.Validated[getPanelsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, installationController.pagination.DefaultPage, installationController.pagination.DefaultPageSize)

	listInfo := installationdto.AdminInstallationListRequest{
		Status: params.Status,
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}
	requests, count, err := installationController.installationService.GetPanelsByAdmin(listInfo)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(requests, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (installationController *AdminInstallationController) SearchPanels(ctx *gin.Context) {
	type searchPanelsParams struct {
		Query    string `form:"query" validate:"required"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		SortBy   uint   `form:"sortBy"`
		Asc      bool   `form:"asc"`
	}

	params := controller.Validated[searchPanelsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, installationController.pagination.DefaultPage, installationController.pagination.DefaultPageSize)

	request := installationdto.SearchPanelsRequest{
		Query:  params.Query,
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}

	panels, count, err := installationController.installationService.SearchPanels(request)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(panels, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (installationController *AdminInstallationController) GetAllPanelStatuses(ctx *gin.Context) {
	statuses := installationController.installationService.GetPanelStatus()
	controller.Response(ctx, 200, "", statuses)
}

func (installationController *AdminInstallationController) GetPanel(ctx *gin.Context) {
	type installationRequestParams struct {
		PanelID uint `uri:"panelID" validate:"required"`
	}
	params := controller.Validated[installationRequestParams](ctx)

	installationRequest, err := installationController.installationService.GetPanelByAdmin(params.PanelID)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", installationRequest)
}

func (installationController *AdminInstallationController) UpdatePanel(ctx *gin.Context) {
	type installationRequestParams struct {
		PanelID              uint    `uri:"panelID" validate:"required"`
		Name                 *string `json:"name"`
		Status               *uint   `json:"status"`
		BuildingType         *uint   `json:"buildingType"`
		Area                 *uint   `json:"area"`
		Power                *uint   `json:"power"`
		Tilt                 *uint   `json:"tilt"`
		Azimuth              *uint   `json:"azimuth"`
		TotalNumberOfModules *uint   `json:"totalNumberOfModules"`
	}
	params := controller.Validated[installationRequestParams](ctx)

	requestInfo := installationdto.UpdatePanelRequest{
		PanelID:              params.PanelID,
		Name:                 params.Name,
		Status:               params.Status,
		BuildingType:         params.BuildingType,
		Area:                 params.Area,
		Power:                params.Power,
		Tilt:                 params.Tilt,
		Azimuth:              params.Azimuth,
		TotalNumberOfModules: params.TotalNumberOfModules,
	}
	if err := installationController.installationService.UpdatePanel(requestInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updatePanel")
	controller.Response(ctx, 201, message, nil)
}

func (installationController *AdminInstallationController) DeletePanel(ctx *gin.Context) {
	type installationRequestParams struct {
		PanelID uint `uri:"panelID" validate:"required"`
	}
	params := controller.Validated[installationRequestParams](ctx)

	if err := installationController.installationService.DeletePanel(params.PanelID); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deletePanel")
	controller.Response(ctx, 200, message, nil)
}
