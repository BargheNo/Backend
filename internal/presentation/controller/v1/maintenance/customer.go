package maintenance

import (
	"strconv"

	"github.com/BargheNo/Backend/bootstrap"
	maintenancedto "github.com/BargheNo/Backend/internal/application/dto/maintenance"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerMaintenanceController struct {
	constants          *bootstrap.Constants
	pagination         *bootstrap.Pagination
	maintenanceService service.MaintenanceService
}

func NewCustomerMaintenanceController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	maintenanceService service.MaintenanceService,
) *CustomerMaintenanceController {
	return &CustomerMaintenanceController{
		constants:          constants,
		pagination:         pagination,
		maintenanceService: maintenanceService,
	}
}

func (maintenanceController *CustomerMaintenanceController) CreateMaintenanceRequest(ctx *gin.Context) {
	type maintenanceRequestParams struct {
		PanelID       uint   `json:"panelID" validate:"required"`
		CorporationID uint   `uri:"corporationID" validate:"required"`
		Subject       string `json:"subject" validate:"required"`
		Description   string `json:"description" validate:"required"`
		UrgencyLevel  uint   `json:"urgencyLevel" validate:"required"`
	}
	params := controller.Validated[maintenanceRequestParams](ctx)
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)
	requestInfo := maintenancedto.NewMaintenanceRequest{
		PanelID:       params.PanelID,
		OwnerID:       ownerID.(uint),
		CorporationID: params.CorporationID,
		Subject:       params.Subject,
		Description:   params.Description,
		UrgencyLevel:  enum.UrgencyLevel(params.UrgencyLevel),
	}

	maintenanceController.maintenanceService.CreateMaintenanceRequest(requestInfo)

	trans := controller.GetTranslator(ctx, maintenanceController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.maintenanceRequest")
	controller.Response(ctx, 201, message, nil)
}

func (maintenanceController *CustomerMaintenanceController) GetCustomerMaintenanceRequests(ctx *gin.Context) {
	ownerID, _ := ctx.Get(maintenanceController.constants.Context.ID)
	defaultPage, err := strconv.Atoi(maintenanceController.pagination.DefaultPage)
	if err != nil {
		defaultPage = 1
	}
	defaultPageSize, err := strconv.Atoi(maintenanceController.pagination.DefaultPageSize)
	if err != nil {
		defaultPageSize = 10
	}
	params := controller.GetPagination(ctx, defaultPage, defaultPageSize)
	offset, limit := params.GetOffsetLimit()
	listInfo := maintenancedto.MaintenanceListRequest{
		OwnerID: ownerID.(uint),
		Offset:  offset,
		Limit:   limit,
	}
	requests := maintenanceController.maintenanceService.GetCustomerMaintenanceRequests(listInfo)
	controller.Response(ctx, 200, "", requests)
}
