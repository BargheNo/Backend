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

type CustomerInstallationController struct {
	constants           *bootstrap.Constants
	pagination          *bootstrap.Pagination
	installationService service.InstallationService
}

func NewCustomerInstallationController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	installationService service.InstallationService,
) *CustomerInstallationController {
	return &CustomerInstallationController{
		constants:           constants,
		pagination:          pagination,
		installationService: installationService,
	}
}

func (installationController *CustomerInstallationController) CreateInstallationRequest(ctx *gin.Context) {
	type installationRequestParams struct {
		Name          string  `json:"name" validate:"required"`
		Area          uint    `json:"area"`
		Power         uint    `json:"power" validate:"required"`
		MaxCost       float64 `json:"maxCost"`
		BuildingType  string  `json:"buildingType" validate:"required"`
		Description   string  `json:"description"`
		Province      string  `json:"province" validate:"required"`
		City          string  `json:"city" validate:"required"`
		StreetAddress string  `json:"streetAddress" validate:"required"`
		PostalCode    string  `json:"postalCode" validate:"required"`
		HouseNumber   string  `json:"houseNumber" validate:"required"`
		Unit          uint    `json:"unit" validate:"required"`
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
		Address: addressdto.CreateAddressRequest{
			Province:      params.Province,
			City:          params.City,
			StreetAddress: params.StreetAddress,
			PostalCode:    params.PostalCode,
			HouseNumber:   params.HouseNumber,
			Unit:          params.Unit,
			OwnerID:       ownerID.(uint),
			OwnerType:     "installation_requests",
		},
	}
	installationController.installationService.CreateInstallationRequest(requestInfo)

	trans := controller.GetTranslator(ctx, installationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.installationRequest")
	controller.Response(ctx, 200, message, nil)
}

func (installationController *CustomerInstallationController) GetOwnerInstallationRequests(ctx *gin.Context) {
	ownerID, _ := ctx.Get(installationController.constants.Context.ID)
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
	listInfo := installationdto.ListOwnerRequestsRequest{
		OwnerID: ownerID.(uint),
		Offset:  offset,
		Limit:   limit,
	}
	requests := installationController.installationService.GetOwnerInstallationRequests(listInfo)
	controller.Response(ctx, 200, "", requests)
}
