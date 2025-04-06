package installation

import (
	"strconv"

	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
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

func (installationController *CorporationInstallationController) AddPanel(ctx *gin.Context) {
	type addPanelParams struct {
		CorporationID        uint   `uri:"corporationID" validate:"required"`
		PanelName            string `json:"panelName" validate:"required"`
		CustomerPhone        string `json:"customerPhone" validate:"required"`
		Power                uint   `json:"power" validate:"required"`
		Area                 uint   `json:"area" validate:"required"`
		BuildingType         string `json:"buildingType" validate:"required"`
		Tilt                 uint   `json:"tilt" validate:"required"`
		Azimuth              uint   `json:"azimuth" validate:"required"`
		TotalNumberOfModules uint   `json:"totalNumberOfModules" validate:"required"`
		ProvinceID           uint   `json:"provinceID" validate:"required"`
		CityID               uint   `json:"cityID" validate:"required"`
		StreetAddress        string `json:"streetAddress" validate:"required"`
		PostalCode           string `json:"postalCode" validate:"required"`
		HouseNumber          string `json:"houseNumber" validate:"required"`
		Unit                 uint   `json:"unit" validate:"required"`
	}
	params := controller.Validated[addPanelParams](ctx)
	operatorID, _ := ctx.Get(installationController.constants.Context.ID)

	panelInfo := installationdto.AddPanelRequest{
		CorporationID:        params.CorporationID,
		OperatorID:           operatorID.(uint),
		CustomerPhone:        params.CustomerPhone,
		Power:                params.Power,
		Area:                 params.Area,
		BuildingType:         params.BuildingType,
		Tilt:                 params.Tilt,
		Azimuth:              params.Azimuth,
		TotalNumberOfModules: params.TotalNumberOfModules,
		Address: addressdto.CreateAddressRequest{
			ProvinceID:    params.ProvinceID,
			CityID:        params.CityID,
			StreetAddress: params.StreetAddress,
			PostalCode:    params.PostalCode,
			HouseNumber:   params.HouseNumber,
			Unit:          params.Unit,
		},
	}
	installationController.installationService.AddPanel(panelInfo)

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.addPanel")
	controller.Response(ctx, 200, message, nil)
}

func (installationController *CorporationInstallationController) GetCorporationPanels(ctx *gin.Context) {
	type getPanelParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[getPanelParams](ctx)
	defaultPage, err := strconv.Atoi(installationController.pagination.DefaultPage)
	if err != nil {
		defaultPage = 1
	}
	defaultPageSize, err := strconv.Atoi(installationController.pagination.DefaultPageSize)
	if err != nil {
		defaultPageSize = 10
	}
	pagination := controller.GetPagination(ctx, defaultPage, defaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	listInfo := installationdto.PanelListRequest{
		CorporationID: params.CorporationID,
		Offset:        offset,
		Limit:         limit,
	}
	panels := installationController.installationService.GetCorporationPanels(listInfo)

	controller.Response(ctx, 200, "", panels)
}
