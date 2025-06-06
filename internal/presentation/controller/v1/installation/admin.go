package installation

import (
	"github.com/BargheNo/Backend/bootstrap"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminInstallationController struct {
	constants           *bootstrap.Constants
	pagination          *bootstrap.Pagination
	installationService service.InstallationService
}

func NewAdminInstallationController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	installationService service.InstallationService,
) *AdminInstallationController {
	return &AdminInstallationController{
		constants:           constants,
		pagination:          pagination,
		installationService: installationService,
	}
}

func (installationController *AdminInstallationController) GetInstallationRequests(ctx *gin.Context) {
	type getRequestsParams struct {
		Status uint `form:"status" validate:"required"`
	}
	params := controller.Validated[getRequestsParams](ctx)

	pagination := controller.GetPagination(ctx, installationController.pagination.DefaultPage, installationController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	listInfo := installationdto.AdminRequestsListRequest{
		Status: params.Status,
		Offset: offset,
		Limit:  limit,
	}
	requests := installationController.installationService.GetInstallationRequestsByAdmin(listInfo)

	controller.Response(ctx, 200, "", requests)
}

func (installationController *AdminInstallationController) GetInstallationRequest(ctx *gin.Context) {
	type installationRequestParams struct {
		RequestID uint `uri:"requestID" validate:"required"`
	}
	params := controller.Validated[installationRequestParams](ctx)

	installationRequest := installationController.installationService.GetPublicInstallationRequest(params.RequestID)

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
	installationController.installationService.UpdateInstallationRequestByAdmin(requestInfo)

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateInstallationRequest")
	controller.Response(ctx, 201, message, nil)
}

func (installationController *AdminInstallationController) DeleteInstallationRequest(ctx *gin.Context) {
	type installationRequestParams struct {
		RequestID uint `uri:"requestID" validate:"required"`
	}
	params := controller.Validated[installationRequestParams](ctx)

	installationController.installationService.DeleteInstallationRequest(params.RequestID)

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteInstallationRequest")
	controller.Response(ctx, 200, message, nil)
}
