package service

import (
	chatdto "github.com/BargheNo/Backend/internal/application/dto/chat"
)

type ChatService interface {
	CreateOrGetRoom(chatdto.CreateOrGetRoomRequest) chatdto.ChatRoomDetailsResponse
	GetUserRooms(chatdto.GetRoomMessageRequest) []chatdto.ChatRoomDetailsResponse
	SaveMessage(roomID, senderID uint, content string) error
	GetRoomMessages(roomID uint) []chatdto.RoomMessagesResponse
	AddParticipant(roomID, userID uint) error
	RemoveParticipant(roomID, userID uint) error
}
