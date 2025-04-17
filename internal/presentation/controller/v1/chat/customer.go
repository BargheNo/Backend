package chat

import (
	"github.com/BargheNo/Backend/bootstrap"
	chatdto "github.com/BargheNo/Backend/internal/application/dto/chat"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/infrastructure/websocket"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerChatController struct {
	constants        *bootstrap.Constants
	pagination       *bootstrap.Pagination
	websocketSetting *bootstrap.WebsocketSetting
	chatService      service.ChatService
	jwtService       service.JWTService
	hub              *websocket.Hub
}

func NewCustomerChatController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	websocketSetting *bootstrap.WebsocketSetting,
	chatService service.ChatService,
	jwtService service.JWTService,
	hub *websocket.Hub,
) *CustomerChatController {
	return &CustomerChatController{
		constants:        constants,
		pagination:       pagination,
		websocketSetting: websocketSetting,
		chatService:      chatService,
		jwtService:       jwtService,
		hub:              hub,
	}
}

func (chatController *CustomerChatController) CreateOrGetRoom(ctx *gin.Context) {
	type roomParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[roomParams](ctx)
	userID, _ := ctx.Get(chatController.constants.Context.ID)

	roomInfo := chatdto.CreateOrGetUserRoomRequest{
		CorporationID: params.CorporationID,
		UserID:        userID.(uint),
	}
	roomsDetails := chatController.chatService.CreateOrGetRoom(roomInfo)

	controller.Response(ctx, 200, "", roomsDetails)
}

func (chatController *CustomerChatController) GetUserRooms(ctx *gin.Context) {
	userID, _ := ctx.Get(chatController.constants.Context.ID)
	roomsDetails := chatController.chatService.GetUserRooms(userID.(uint))
	controller.Response(ctx, 200, "", roomsDetails)
}

func (chatController *CustomerChatController) GetMessages(ctx *gin.Context) {
	type getMessagesParams struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	userID, _ := ctx.Get(chatController.constants.Context.ID)
	pagination := controller.GetPagination(ctx, chatController.pagination.DefaultPage, chatController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()
	param := controller.Validated[getMessagesParams](ctx)

	roomInfo := chatdto.GetRoomMessageRequest{
		RoomID: param.RoomID,
		UserID: userID.(uint),
		Offset: offset,
		Limit:  limit,
	}
	messages := chatController.chatService.GetRoomMessages(roomInfo)

	controller.Response(ctx, 200, "", messages)
}

func (chatController *CustomerChatController) HandleWebsocket(ctx *gin.Context) {
	type roomConnectionParams struct {
		RoomID uint   `uri:"roomID" validate:"required"`
		Token  string `uri:"token" validate:"required"`
	}
	param := controller.Validated[roomConnectionParams](ctx)

	claims, err := chatController.jwtService.ValidateToken(param.Token)
	if err != nil {
		panic(err)
	}
	userID := uint(claims["sub"].(float64))
	conn, _ := ctx.Get(chatController.constants.Context.WebsocketConnection)

	client := websocket.NewClient(chatController.hub, conn, param.RoomID, userID, chatController.websocketSetting, chatController.chatService, nil)
	client.Hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
