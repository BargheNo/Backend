package monitoring

import (
	"github.com/BargheNo/Backend/bootstrap"
	monitoringdto "github.com/BargheNo/Backend/internal/application/dto/monitoring"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/infrastructure/websocket"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerMonitoringController struct {
	constants         *bootstrap.Constants
	hub               *websocket.Hub
	jwtService        usecase.JWTService
	websocketSetting  *bootstrap.WebsocketSetting
	monitoringService usecase.MonitoringService
	pagination        *bootstrap.Pagination
}

func NewCustomerMonitoringController(
	constants *bootstrap.Constants,
	hub *websocket.Hub,
	jwtService usecase.JWTService,
	websocketSetting *bootstrap.WebsocketSetting,
	monitoringService usecase.MonitoringService,
	pagination *bootstrap.Pagination,
) *CustomerMonitoringController {
	return &CustomerMonitoringController{
		constants:         constants,
		hub:               hub,
		jwtService:        jwtService,
		websocketSetting:  websocketSetting,
		monitoringService: monitoringService,
		pagination:        pagination,
	}
}

func (monitoringController *CustomerMonitoringController) HandleWebsocket(ctx *gin.Context) {
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

	client := websocket.NewClient(
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

func (monitoringController *CustomerMonitoringController) GetPanelStatus(ctx *gin.Context) {
	type getPanelStatusParams struct {
		PanelID  uint `uri:"panelID" validate:"required"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
	}
	param := controller.Validated[getPanelStatusParams](ctx)

	ownerID, _ := ctx.Get(monitoringController.constants.Context.ID)

	offset, limit := controller.GetOffsetLimit(param.Page, param.PageSize, monitoringController.pagination.DefaultPage, monitoringController.pagination.DefaultPageSize)

	listInfo := monitoringdto.CustomerPanelStatusListRequest{
		PanelID: param.PanelID,
		OwnerID: ownerID.(uint),
		Offset:  offset,
		Limit:   limit,
	}

	response, count, err := monitoringController.monitoringService.GetPanelStatus(listInfo)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(response, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (monitoringController *CustomerMonitoringController) GetPanelHistory(ctx *gin.Context) {
	type getPanelHistoryParams struct {
		PanelID  uint `uri:"panelID" validate:"required"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
	}
	param := controller.Validated[getPanelHistoryParams](ctx)

	ownerID, _ := ctx.Get(monitoringController.constants.Context.ID)

	offset, limit := controller.GetOffsetLimit(param.Page, param.PageSize, monitoringController.pagination.DefaultPage, monitoringController.pagination.DefaultPageSize)

	listInfo := monitoringdto.CustomerPanelStatusListRequest{
		PanelID: param.PanelID,
		OwnerID: ownerID.(uint),
		Offset:  offset,
		Limit:   limit,
	}

	response, count, err := monitoringController.monitoringService.GetPanelHistory(listInfo)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(response, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (monitoringController *CustomerMonitoringController) GetPanelEvent(ctx *gin.Context) {
	type getPanelEventParams struct {
		PanelID  uint `uri:"panelID" validate:"required"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
	}
	param := controller.Validated[getPanelEventParams](ctx)

	ownerID, _ := ctx.Get(monitoringController.constants.Context.ID)

	offset, limit := controller.GetOffsetLimit(param.Page, param.PageSize, monitoringController.pagination.DefaultPage, monitoringController.pagination.DefaultPageSize)

	listInfo := monitoringdto.CustomerPanelStatusListRequest{
		PanelID: param.PanelID,
		OwnerID: ownerID.(uint),
		Offset:  offset,
		Limit:   limit,
	}

	response, count, err := monitoringController.monitoringService.GetPanelEvent(listInfo)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(response, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}
