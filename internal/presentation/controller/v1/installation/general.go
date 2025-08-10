package installation

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralInstallationController struct {
	constants           *bootstrap.Constants
	installationService usecase.InstallationService
}

func NewGeneralInstallationController(
	constants *bootstrap.Constants,
	installationService usecase.InstallationService,
) *GeneralInstallationController {
	return &GeneralInstallationController{
		constants:           constants,
		installationService: installationService,
	}
}

func (installationController *GeneralInstallationController) GetRequestStatuses(ctx *gin.Context) {
	statuses := installationController.installationService.GetRequestStatuses()
	controller.Response(ctx, 200, "", statuses)
}

func (installationController *GeneralInstallationController) GetPanelStatuses(ctx *gin.Context) {
	statuses := installationController.installationService.GetPanelStatuses()
	controller.Response(ctx, 200, "", statuses)
}

func (installationController *GeneralInstallationController) GetBuildingTypes(ctx *gin.Context) {
	types := installationController.installationService.GetBuildingTypes()
	controller.Response(ctx, 200, "", types)
}

func (installationController *GeneralInstallationController) GetRequestSortableFields(ctx *gin.Context) {
	columns := installationController.installationService.GetRequestSortableColumns()
	controller.Response(ctx, 200, "", columns)
}

func (installationController *GeneralInstallationController) GetPanelSortableFields(ctx *gin.Context) {
	columns := installationController.installationService.GetPanelSortableColumns()
	controller.Response(ctx, 200, "", columns)
}
