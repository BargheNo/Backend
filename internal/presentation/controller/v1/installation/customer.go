package installation

import (
	"github.com/BargheNo/Backend/bootstrap"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
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

func (installationController *CustomerInstallationController) CreateInstallationRequest(ctx *gin.Context) {
	type installationRequestParams struct {
		Name         string  `json:"name" validate:"required"`
		Area         uint    `json:"area"`
		Power        uint    `json:"power" validate:"required"`
		MaxCost      float64 `json:"maxCost"`
		BuildingType string  `json:"buildingType" validate:"required"`
		Description  string  `json:"description"`
		AddressID    uint    `json:"addressID" validate:"required"`
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
		AddressID:    params.AddressID,
	}
	installationController.installationService.CreateInstallationRequest(requestInfo)

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.installationRequest")
	controller.Response(ctx, 200, message, nil)
}

func (installationController *CustomerInstallationController) GetInstallationRequests(ctx *gin.Context) {
	// type installationRequestsListParams struct {
	// 	Description string `json:"description"`
	// 	AddressID   uint   `json:"addressID" validate:"required"`
	// }
	// params := controller.Validated[installationRequestsListParams](ctx)
}
