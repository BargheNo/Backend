package monitoring

import (
	"github.com/BargheNo/Backend/bootstrap"
	monitoringdto "github.com/BargheNo/Backend/internal/application/dto/monitoring"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CorporationMonitoringController struct {
	constants         *bootstrap.Constants
	monitoringService usecase.MonitoringService
	pagination        *bootstrap.Pagination
}

func NewCorporationMonitoringController(
	constants *bootstrap.Constants,
	monitoringService usecase.MonitoringService,
	pagination *bootstrap.Pagination,
) *CorporationMonitoringController {
	return &CorporationMonitoringController{
		constants:         constants,
		monitoringService: monitoringService,
		pagination:        pagination,
	}
}

func (monitoringController *CorporationMonitoringController) GetPanelStatus(ctx *gin.Context) {
	type getPanelStatusParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		PanelID       uint `uri:"panelID" validate:"required"`
		Page          int  `form:"page"`
		PageSize      int  `form:"pageSize"`
	}
	param := controller.Validated[getPanelStatusParams](ctx)

	offset, limit := controller.GetOffsetLimit(param.Page, param.PageSize, monitoringController.pagination.DefaultPage, monitoringController.pagination.DefaultPageSize)

	userID, _ := ctx.Get(monitoringController.constants.Context.ID)

	listInfo := monitoringdto.CorporationPanelStatusListRequest{
		PanelID:       param.PanelID,
		CorporationID: param.CorporationID,
		UserID:        userID.(uint),
		Offset:        offset,
		Limit:         limit,
	}

	response, count, err := monitoringController.monitoringService.GetCorporationPanelStatus(listInfo)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(response, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (monitoringController *CorporationMonitoringController) GetPanelHistory(ctx *gin.Context) {
	type getPanelHistoryParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		PanelID       uint `uri:"panelID" validate:"required"`
		Page          int  `form:"page"`
		PageSize      int  `form:"pageSize"`
	}
	param := controller.Validated[getPanelHistoryParams](ctx)

	offset, limit := controller.GetOffsetLimit(param.Page, param.PageSize, monitoringController.pagination.DefaultPage, monitoringController.pagination.DefaultPageSize)

	userID, _ := ctx.Get(monitoringController.constants.Context.ID)

	listInfo := monitoringdto.CorporationPanelStatusListRequest{
		PanelID:       param.PanelID,
		CorporationID: param.CorporationID,
		UserID:        userID.(uint),
		Offset:        offset,
		Limit:         limit,
	}

	response, count, err := monitoringController.monitoringService.GetCorporationPanelHistory(listInfo)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(response, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (monitoringController *CorporationMonitoringController) GetPanelEvent(ctx *gin.Context) {
	type getPanelEventParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		PanelID       uint `uri:"panelID" validate:"required"`
		Page          int  `form:"page"`
		PageSize      int  `form:"pageSize"`
	}
	param := controller.Validated[getPanelEventParams](ctx)

	offset, limit := controller.GetOffsetLimit(param.Page, param.PageSize, monitoringController.pagination.DefaultPage, monitoringController.pagination.DefaultPageSize)

	userID, _ := ctx.Get(monitoringController.constants.Context.ID)

	listInfo := monitoringdto.CorporationPanelStatusListRequest{
		PanelID:       param.PanelID,
		CorporationID: param.CorporationID,
		UserID:        userID.(uint),
		Offset:        offset,
		Limit:         limit,
	}

	response, count, err := monitoringController.monitoringService.GetCorporationPanelEvent(listInfo)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(response, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}
