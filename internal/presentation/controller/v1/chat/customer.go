package chat

import (
	"github.com/BargheNo/Backend/bootstrap"
	chatdto "github.com/BargheNo/Backend/internal/application/dto/chat"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/websocket"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerChatController struct {
	constants        *bootstrap.Constants
	pagination       *bootstrap.Pagination
	websocketSetting *bootstrap.WebsocketSetting
	chatService      usecase.ChatService
	jwtService       usecase.JWTService
	userService      usecase.UserService
	hub              *websocket.Hub
}

func NewCustomerChatController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	websocketSetting *bootstrap.WebsocketSetting,
	chatService usecase.ChatService,
	jwtService usecase.JWTService,
	userService usecase.UserService,
	hub *websocket.Hub,
) *CustomerChatController {
	return &CustomerChatController{
		constants:        constants,
		pagination:       pagination,
		websocketSetting: websocketSetting,
		chatService:      chatService,
		jwtService:       jwtService,
		userService:      userService,
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
	roomsDetails, err := chatController.chatService.CreateOrGetRoom(roomInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", roomsDetails)
}

func (chatController *CustomerChatController) GetUserRooms(ctx *gin.Context) {
	userID, _ := ctx.Get(chatController.constants.Context.ID)
	roomsDetails, err := chatController.chatService.GetUserRooms(userID.(uint))
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", roomsDetails)
}

func (chatController *CustomerChatController) GetMessages(ctx *gin.Context) {
	type getMessagesParams struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	param := controller.Validated[getMessagesParams](ctx)
	userID, _ := ctx.Get(chatController.constants.Context.ID)
	pagination := controller.GetPagination(ctx, chatController.pagination.DefaultPage, chatController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	roomInfo := chatdto.GetRoomMessageRequest{
		RoomID: param.RoomID,
		UserID: userID.(uint),
		Offset: offset,
		Limit:  limit,
	}
	messages, err := chatController.chatService.GetRoomMessages(roomInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", messages)
}

func (chatController *CustomerChatController) BlockRoom(ctx *gin.Context) {
	type getMessagesParams struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	param := controller.Validated[getMessagesParams](ctx)
	userID, _ := ctx.Get(chatController.constants.Context.ID)

	blockRequest := chatdto.BlockServiceChatRequest{
		UserID:     userID.(uint),
		RoomID:     param.RoomID,
		BlockedBy:  enum.BlockedByUser,
		ChatStatus: enum.ChatStatusBlocked,
	}
	if err := chatController.chatService.BlockChatRoom(blockRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, chatController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.blockChatRoom")
	controller.Response(ctx, 200, message, nil)
}

func (chatController *CustomerChatController) UnBlockRoom(ctx *gin.Context) {
	type getMessagesParams struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	param := controller.Validated[getMessagesParams](ctx)
	userID, _ := ctx.Get(chatController.constants.Context.ID)

	blockRequest := chatdto.BlockServiceChatRequest{
		UserID:     userID.(uint),
		RoomID:     param.RoomID,
		BlockedBy:  enum.BlockedByUser,
		ChatStatus: enum.ChatStatusActive,
	}
	if err := chatController.chatService.UnBlockChatRoom(blockRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, chatController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.unblockChatRoom")
	controller.Response(ctx, 200, message, nil)
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
