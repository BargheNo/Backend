package notification

import (
	"github.com/BargheNo/Backend/bootstrap"
	notificationdto "github.com/BargheNo/Backend/internal/application/dto/notification"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/infrastructure/websocket"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerNotificationController struct {
	constants           *bootstrap.Constants
	websocketSetting    *bootstrap.WebsocketSetting
	pagination          *bootstrap.Pagination
	notificationService usecase.NotificationService
	jwtService          usecase.JWTService
	userService         usecase.UserService
	hub                 *websocket.Hub
}

func NewCustomerNotificationController(
	constants *bootstrap.Constants,
	websocketSetting *bootstrap.WebsocketSetting,
	pagination *bootstrap.Pagination,
	notificationService usecase.NotificationService,
	jwtService usecase.JWTService,
	userService usecase.UserService,
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
	if err := notificationController.notificationService.MarkAsRead(notificationInfo); err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "successMessage.readNotification", nil)
}

func (notificationController *CustomerNotificationController) GetUserNotifications(ctx *gin.Context) {
	type notificationsParams struct {
		Types    []uint `form:"notificationTypes" validate:"required"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		SortBy   uint   `form:"sortBy"`
		Asc      bool   `form:"asc"`
	}
	params := controller.Validated[notificationsParams](ctx)
	userID, _ := ctx.Get(notificationController.constants.Context.ID)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, notificationController.pagination.DefaultPage, notificationController.pagination.DefaultPageSize)

	notificationsRequest := notificationdto.NotificationListRequest{
		Types:  params.Types,
		UserID: userID.(uint),
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}
	notificationsDetails, count, err := notificationController.notificationService.GetUserNotifications(notificationsRequest)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(notificationsDetails, count, params.Page, params.PageSize)

	controller.Response(ctx, 200, "", data)
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
	if err := notificationController.notificationService.UpdateNotificationSettings(settingInfo); err != nil {
		panic(err)
	}

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
