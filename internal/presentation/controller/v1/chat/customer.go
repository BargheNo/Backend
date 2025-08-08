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
		RoomID   uint `uri:"roomID" validate:"required"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
		SortBy   uint `form:"sortBy"`
		Asc      bool `form:"asc"`
	}
	params := controller.Validated[getMessagesParams](ctx)
	userID, _ := ctx.Get(chatController.constants.Context.ID)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, chatController.pagination.DefaultPage, chatController.pagination.DefaultPageSize)

	roomInfo := chatdto.GetRoomMessageRequest{
		RoomID: params.RoomID,
		UserID: userID.(uint),
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}
	messages, count, err := chatController.chatService.GetRoomMessages(roomInfo)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(messages, count, params.Page, params.PageSize)

	controller.Response(ctx, 200, "", data)
}

func (chatController *CustomerChatController) BlockRoom(ctx *gin.Context) {
	type getMessagesParams struct {
		RoomID uint `uri:"roomID" validate:"required"`
	}
	params := controller.Validated[getMessagesParams](ctx)
	userID, _ := ctx.Get(chatController.constants.Context.ID)

	blockRequest := chatdto.BlockServiceChatRequest{
		UserID:     userID.(uint),
		RoomID:     params.RoomID,
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
	params := controller.Validated[getMessagesParams](ctx)
	userID, _ := ctx.Get(chatController.constants.Context.ID)

	blockRequest := chatdto.BlockServiceChatRequest{
		UserID:     userID.(uint),
		RoomID:     params.RoomID,
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
	params := controller.Validated[roomConnectionParams](ctx)

	claims, err := chatController.jwtService.ValidateToken(params.Token)
	if err != nil {
		panic(err)
	}
	userID := uint(claims["sub"].(float64))
	conn, _ := ctx.Get(chatController.constants.Context.WebsocketConnection)

	client := websocket.NewClient(chatController.hub, conn, params.RoomID, userID, chatController.websocketSetting, chatController.chatService, nil)
	client.Hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
