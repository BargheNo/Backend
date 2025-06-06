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
