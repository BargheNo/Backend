package service

import (
	chatdto "github.com/BargheNo/Backend/internal/application/dto/chat"
)

type ChatService interface {
	CreateOrGetRoom(request chatdto.CreateOrGetRoomRequest) chatdto.ChatRoomDetailsResponse
	GetUserRooms(userID uint) []chatdto.ChatRoomDetailsResponse
	SaveMessage(roomID, senderID uint, content string)
	GetRoomMessages(request chatdto.GetRoomMessageRequest) []chatdto.RoomMessagesResponse
}
