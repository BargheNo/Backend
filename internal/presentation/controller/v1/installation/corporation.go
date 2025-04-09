package installation

import (
	"github.com/BargheNo/Backend/bootstrap"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CorporationInstallationController struct {
	constants           *bootstrap.Constants
	pagination          *bootstrap.Pagination
	installationService service.InstallationService
}

func NewCorporationInstallationController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	installationService service.InstallationService,
) *CorporationInstallationController {
	return &CorporationInstallationController{
		constants:           constants,
		pagination:          pagination,
		installationService: installationService,
	}
}

func (installationController *CorporationInstallationController) GetInstallationRequests(ctx *gin.Context) {
	// refactor to support status
	corporationID, _ := ctx.Get(installationController.constants.Context.ID)
	params := controller.GetPagination(ctx, installationController.pagination.DefaultPage, installationController.pagination.DefaultPageSize)
	offset, limit := params.GetOffsetLimit()
	listInfo := installationdto.InstallationListRequest{
		OwnerID: corporationID.(uint),
		Offset:  offset,
		Limit:   limit,
	}
	installationRequest := installationController.installationService.GetInstallationRequests(listInfo)

	controller.Response(ctx, 200, "", installationRequest)
}
