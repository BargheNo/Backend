package monitoring

import (
	"github.com/BargheNo/Backend/bootstrap"
	monitoringdto "github.com/BargheNo/Backend/internal/application/dto/monitoring"
	"github.com/BargheNo/Backend/internal/application/usecase"
	ws "github.com/BargheNo/Backend/internal/infrastructure/websocket"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminMonitoringController struct {
	constants         *bootstrap.Constants
	monitoringService usecase.MonitoringService
	pagination        *bootstrap.Pagination
	hub               *ws.Hub
	jwtService        usecase.JWTService
	websocketSetting  *bootstrap.WebsocketSetting
}

func NewAdminMonitoringController(
	constants *bootstrap.Constants,
	monitoringService usecase.MonitoringService,
	pagination *bootstrap.Pagination,
	hub *ws.Hub,
	jwtService usecase.JWTService,
	websocketSetting *bootstrap.WebsocketSetting,
) *AdminMonitoringController {
	return &AdminMonitoringController{
		constants:         constants,
		monitoringService: monitoringService,
		pagination:        pagination,
		hub:               hub,
		jwtService:        jwtService,
		websocketSetting:  websocketSetting,
	}
}

func (monitoringController *AdminMonitoringController) HandleWebsocket(ctx *gin.Context) {
	type roomConnectionParams struct {
		PanelID uint   `uri:"panelID" validate:"required"`
		Token   string `uri:"token" validate:"required"`
	}
	param := controller.Validated[roomConnectionParams](ctx)

	claims, err := monitoringController.jwtService.ValidateToken(param.Token)
	if err != nil {
		panic(err)
	}
	userID := uint(claims["sub"].(float64))

	conn, _ := ctx.Get(monitoringController.constants.Context.WebsocketConnection)

	client := ws.NewClient(
		monitoringController.hub,
		conn,
		param.PanelID,
		userID,
		monitoringController.websocketSetting,
		nil,
		nil,
	)
	client.Hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
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
