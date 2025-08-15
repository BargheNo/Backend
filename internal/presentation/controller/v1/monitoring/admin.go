package monitoring

import (
	"github.com/BargheNo/Backend/bootstrap"
	monitoringdto "github.com/BargheNo/Backend/internal/application/dto/monitoring"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminMonitoringController struct {
	monitoringService usecase.MonitoringService
	pagination        *bootstrap.Pagination
}

func NewAdminMonitoringController(
	monitoringService usecase.MonitoringService,
	pagination *bootstrap.Pagination,
) *AdminMonitoringController {
	return &AdminMonitoringController{
		monitoringService: monitoringService,
		pagination:        pagination,
	}
}

func (monitoringController *AdminMonitoringController) GetPanelStatus(ctx *gin.Context) {
	type getPanelStatusParams struct {
		PanelID  uint `uri:"panelID" validate:"required"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
	}
	param := controller.Validated[getPanelStatusParams](ctx)

	offset, limit := controller.GetOffsetLimit(param.Page, param.PageSize, monitoringController.pagination.DefaultPage, monitoringController.pagination.DefaultPageSize)

	listInfo := monitoringdto.AdminPanelStatusListRequest{
		PanelID: param.PanelID,
		Offset:  offset,
		Limit:   limit,
	}

	response, count, err := monitoringController.monitoringService.GetAdminPanelStatus(listInfo)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(response, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (monitoringController *AdminMonitoringController) GetPanelHistory(ctx *gin.Context) {
	type getPanelHistoryParams struct {
		PanelID  uint `uri:"panelID" validate:"required"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
	}
	param := controller.Validated[getPanelHistoryParams](ctx)

	offset, limit := controller.GetOffsetLimit(param.Page, param.PageSize, monitoringController.pagination.DefaultPage, monitoringController.pagination.DefaultPageSize)

	listInfo := monitoringdto.AdminPanelStatusListRequest{
		PanelID: param.PanelID,
		Offset:  offset,
		Limit:   limit,
	}

	response, count, err := monitoringController.monitoringService.GetAdminPanelHistory(listInfo)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(response, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (monitoringController *AdminMonitoringController) GetPanelEvent(ctx *gin.Context) {
	type getPanelEventParams struct {
		PanelID  uint `uri:"panelID" validate:"required"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
	}
	param := controller.Validated[getPanelEventParams](ctx)

	offset, limit := controller.GetOffsetLimit(param.Page, param.PageSize, monitoringController.pagination.DefaultPage, monitoringController.pagination.DefaultPageSize)

	listInfo := monitoringdto.AdminPanelStatusListRequest{
		PanelID: param.PanelID,
		Offset:  offset,
		Limit:   limit,
	}

	response, count, err := monitoringController.monitoringService.GetAdminPanelEvent(listInfo)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(response, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}
