package installation

import (
	"strconv"

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
	defaultPage, err := strconv.Atoi(installationController.pagination.DefaultPage)
	if err != nil {
		defaultPage = 1
	}
	defaultPageSize, err := strconv.Atoi(installationController.pagination.DefaultPageSize)
	if err != nil {
		defaultPageSize = 10
	}
	params := controller.GetPagination(ctx, defaultPage, defaultPageSize)
	offset, limit := params.GetOffsetLimit()
	listInfo := installationdto.InstallationListRequest{
		OwnerID: corporationID.(uint),
		Offset:  offset,
		Limit:   limit,
	}
	installationRequest := installationController.installationService.GetInstallationRequests(listInfo)

	controller.Response(ctx, 200, "", installationRequest)
}
