package monitoring

import (
	"github.com/BargheNo/Backend/bootstrap"
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
}

func NewCustomerMonitoringController(
	constants *bootstrap.Constants,
	hub *websocket.Hub,
	jwtService usecase.JWTService,
	websocketSetting *bootstrap.WebsocketSetting,
	monitoringService usecase.MonitoringService,
) *CustomerMonitoringController {
	return &CustomerMonitoringController{
		constants:         constants,
		hub:               hub,
		jwtService:        jwtService,
		websocketSetting:  websocketSetting,
		monitoringService: monitoringService,
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
