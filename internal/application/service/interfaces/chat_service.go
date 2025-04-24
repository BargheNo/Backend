package service

import (
	chatdto "github.com/BargheNo/Backend/internal/application/dto/chat"
	"github.com/BargheNo/Backend/internal/domain/entity"
)

type ChatService interface {
	CreateChatRoom(request chatdto.CreateOrGetUserRoomRequest) *entity.ChatRoom
	CreateOrGetRoom(request chatdto.CreateOrGetUserRoomRequest) chatdto.ChatRoomDetailsResponse
	GetCorporationRoom(request chatdto.GetCorporationRoomRequest) chatdto.ChatRoomDetailsResponse
	GetUserRooms(userID uint) []chatdto.ChatRoomDetailsResponse
	GetCorporationRooms(request chatdto.GetCorporationRoomsRequest) []chatdto.ChatRoomDetailsResponse
	SaveMessage(roomID, senderID uint, content string) chatdto.RoomMessagesResponse
	GetRoomMessages(request chatdto.GetRoomMessageRequest) []chatdto.RoomMessagesResponse
	BlockChatRoom(request chatdto.BlockServiceChatRequest)
	UnBlockChatRoom(request chatdto.BlockServiceChatRequest)
}
