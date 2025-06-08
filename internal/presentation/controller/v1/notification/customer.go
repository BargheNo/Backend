package notification

import (
	"github.com/BargheNo/Backend/bootstrap"
	notificationdto "github.com/BargheNo/Backend/internal/application/dto/notification"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/infrastructure/websocket"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerNotificationController struct {
	constants           *bootstrap.Constants
	websocketSetting    *bootstrap.WebsocketSetting
	pagination          *bootstrap.Pagination
	notificationService service.NotificationService
	jwtService          service.JWTService
	userService         service.UserService
	hub                 *websocket.Hub
}

func NewCustomerNotificationController(
	constants *bootstrap.Constants,
	websocketSetting *bootstrap.WebsocketSetting,
	pagination *bootstrap.Pagination,
	notificationService service.NotificationService,
	jwtService service.JWTService,
	userService service.UserService,
	hub *websocket.Hub,
) *CustomerNotificationController {
	return &CustomerNotificationController{
		constants:           constants,
		websocketSetting:    websocketSetting,
		pagination:          pagination,
		notificationService: notificationService,
		jwtService:          jwtService,
		userService:         userService,
		hub:                 hub,
	}
}

func (notificationController *CustomerNotificationController) MarkAsRead(ctx *gin.Context) {
	type notificationParams struct {
		NotificationID uint `uri:"notificationID" validate:"required"`
	}
	params := controller.Validated[notificationParams](ctx)
	userID, _ := ctx.Get(notificationController.constants.Context.ID)

	notificationInfo := notificationdto.NotificationInfoRequest{
		NotificationID: params.NotificationID,
		UserID:         userID.(uint),
	}
	notificationController.notificationService.MarkAsRead(notificationInfo)

	controller.Response(ctx, 200, "successMessage.readNotification", nil)
}

func (notificationController *CustomerNotificationController) GetUserNotifications(ctx *gin.Context) {
	type notificationsParams struct {
		Types []uint `form:"notificationTypes" validate:"required"`
	}
	params := controller.Validated[notificationsParams](ctx)
	userID, _ := ctx.Get(notificationController.constants.Context.ID)
	pagination := controller.GetPagination(ctx, notificationController.pagination.DefaultPage, notificationController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	notificationsRequest := notificationdto.NotificationListRequest{
		Types:  params.Types,
		UserID: userID.(uint),
		Offset: offset,
		Limit:  limit,
	}
	notificationsDetails, err := notificationController.notificationService.GetUserNotifications(notificationsRequest)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", notificationsDetails)
}

func (notificationController *CustomerNotificationController) GetUserNotificationSettings(ctx *gin.Context) {
	userID, _ := ctx.Get(notificationController.constants.Context.ID)
	settingsDetails, err := notificationController.notificationService.GetUserNotificationSettings(userID.(uint))
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", settingsDetails)
}

func (notificationController *CustomerNotificationController) UpdateSettings(ctx *gin.Context) {
	type settingsParams struct {
		SettingID      uint `uri:"settingID" validate:"required"`
		IsEmailEnabled bool `json:"isEmailEnabled"`
		IsPushEnabled  bool `json:"isPushEnabled"`
	}
	params := controller.Validated[settingsParams](ctx)
	userID, _ := ctx.Get(notificationController.constants.Context.ID)

	settingInfo := notificationdto.UpdateSettingsRequest{
		SettingID:      params.SettingID,
		UserID:         userID.(uint),
		IsEmailEnabled: params.IsEmailEnabled,
		IsPushEnabled:  params.IsPushEnabled,
	}
	notificationController.notificationService.UpdateNotificationSettings(settingInfo)

	trans := controller.GetTranslator(ctx, notificationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateNotificationSetting")
	controller.Response(ctx, 200, message, nil)
}

func (notificationController *CustomerNotificationController) HandleWebsocket(ctx *gin.Context) {
	type notificationConnectionParams struct {
		Token string `uri:"token" validate:"required"`
	}
	param := controller.Validated[notificationConnectionParams](ctx)

	claims, err := notificationController.jwtService.ValidateToken(param.Token)
	if err != nil {
		panic(err)
	}
	userID := uint(claims["sub"].(float64))
	conn, _ := ctx.Get(notificationController.constants.Context.WebsocketConnection)

	client := websocket.NewClient(notificationController.hub, conn, 0, userID, notificationController.websocketSetting, nil, notificationController.notificationService)
	client.Hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
