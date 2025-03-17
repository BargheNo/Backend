package installation

import (
	"github.com/BargheNo/Backend/bootstrap"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/gin-gonic/gin"
)

type CustomerInstallationController struct {
	constants           *bootstrap.Constants
	installationService service.InstallationService
}

func NewCustomerInstallationController(
	constants *bootstrap.Constants,
	installationService service.InstallationService,
) *CustomerInstallationController {
	return &CustomerInstallationController{
		constants:           constants,
		installationService: installationService,
	}
}

func (installationController *CustomerInstallationController) PanelInstallationRequest(ctx *gin.Context) {
	// type installationRequestParams struct {
	// }
	// params := controller.Validated[installationRequestParams](ctx, &userController.constants.Context)
	// userID, _ := ctx.Get(userController.constants.Context.ID)
	// resetPasswordInfo := customerdto.NewInstallationRequest{}
	// userController.installationService.

	// trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	// message, _ := trans.Translate("successMessage.")
	// controller.Response(ctx, 200, message, nil)
}
