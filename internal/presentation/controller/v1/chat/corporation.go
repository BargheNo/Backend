package chat

import (
	"github.com/BargheNo/Backend/bootstrap"
	chatdto "github.com/BargheNo/Backend/internal/application/dto/chat"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CorporationChatController struct {
	constants   *bootstrap.Constants
	chatService service.ChatService
}

func NewCorporationChatController(
	constants *bootstrap.Constants,
	chatService service.ChatService,
) *CorporationChatController {
	return &CorporationChatController{
		constants:   constants,
		chatService: chatService,
	}
}

func (chatController *CorporationChatController) GetRoom(ctx *gin.Context) {
	type roomParams struct {
		CorporationID uint   `uri:"corporationID" validate:"required"`
		Phone         string `form:"phone" validate:"required,e164"`
	}
	params := controller.Validated[roomParams](ctx)
	userID, _ := ctx.Get(chatController.constants.Context.ID)

	roomInfo := chatdto.GetCorporationRoomRequest{
		CorporationID: params.CorporationID,
		ApplicantID:   userID.(uint),
		UserPhone:     params.Phone,
	}
	roomsDetails := chatController.chatService.GetCorporationRoom(roomInfo)

	controller.Response(ctx, 200, "", roomsDetails)
}

func (chatController *CorporationChatController) GetRooms(ctx *gin.Context) {
	type roomParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[roomParams](ctx)
	userID, _ := ctx.Get(chatController.constants.Context.ID)
	request := chatdto.GetCorporationRoomsRequest{
		CorporationID: params.CorporationID,
		ApplicantID:   userID.(uint),
	}
	roomsDetails := chatController.chatService.GetCorporationRooms(request)
	controller.Response(ctx, 200, "", roomsDetails)
}
