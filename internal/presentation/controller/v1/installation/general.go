package installation

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/application/port"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralInstallationController struct {
	constants           *bootstrap.Constants
	installationService port.InstallationService
}

func NewGeneralInstallationController(
	constants *bootstrap.Constants,
	installationService port.InstallationService,
) *GeneralInstallationController {
	return &GeneralInstallationController{
		constants:           constants,
		installationService: installationService,
	}
}

func (installationController *GeneralInstallationController) GetRequestStatuses(ctx *gin.Context) {
	statuses := installationController.installationService.GetRequestStatuses()
	controller.Response(ctx, 201, "", statuses)
}

func (installationController *GeneralInstallationController) GetBuildingTypes(ctx *gin.Context) {
	types := installationController.installationService.GetBuildingTypes()
	controller.Response(ctx, 201, "", types)
}
