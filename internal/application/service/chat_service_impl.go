package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	chatdto "github.com/BargheNo/Backend/internal/application/dto/chat"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type ChatService struct {
	constants          *bootstrap.Constants
	userService        service.UserService
	corporationService service.CorporationService
	chatRepository     repository.ChatRepository
	db                 database.Database
}

func NewChatService(
	constants *bootstrap.Constants,
	userService service.UserService,
	corporationService service.CorporationService,
	chatRepository repository.ChatRepository,
	db database.Database,
) *ChatService {
	return &ChatService{
		constants:          constants,
		userService:        userService,
		corporationService: corporationService,
		chatRepository:     chatRepository,
		db:                 db,
	}
}

func (chatService *ChatService) CreateOrGetRoom(request chatdto.CreateOrGetRoomRequest) chatdto.ChatRoomDetailsResponse {
	customer := chatService.userService.GetUserCredential(request.UserID)
	corporation := chatService.corporationService.GetCorporationByID(request.CorporationID)
	var room *entity.ChatRoom
	var exist bool
	room, exist = chatService.chatRepository.GetUserAndCorpRoom(chatService.db, request.UserID, request.CorporationID)
	if !exist {
		room = &entity.ChatRoom{
			CorporationID: request.CorporationID,
			CustomerID:    request.UserID,
		}
		chatService.chatRepository.CreateRoom(chatService.db, room)
	}
	roomDetails := chatdto.ChatRoomDetailsResponse{
		RoomID:                room.ID,
		CustomerCredential:    customer,
		CorporationCredential: corporationdto.CorporationDetailsResponse{ID: request.CorporationID, Name: corporation.Name},
	}

	return roomDetails
}

func (chatService *ChatService) GetUserRooms(userID uint) []chatdto.ChatRoomDetailsResponse {
	customer := chatService.userService.GetUserCredential(userID)
	rooms := chatService.chatRepository.GetUserRooms(chatService.db, userID)
	roomsDetails := make([]chatdto.ChatRoomDetailsResponse, len(rooms))
	for i, room := range rooms {
		corporation := chatService.corporationService.GetCorporationByID(room.CorporationID)
		roomsDetails[i] = chatdto.ChatRoomDetailsResponse{
			RoomID:                room.ID,
			CustomerCredential:    customer,
			CorporationCredential: corporationdto.CorporationDetailsResponse{ID: corporation.ID, Name: corporation.Name},
		}
	}
	return roomsDetails
}

func (chatService *ChatService) validateRoomParticipantAccess(senderID, memberID, corporationID uint) {
	if senderID != memberID {
		chatService.corporationService.CheckApplicantAccess(corporationID, senderID)
	}
}

func (chatService *ChatService) SaveMessage(roomID, senderID uint, content string) {
	chatService.userService.GetUserCredential(senderID)
	room, exist := chatService.chatRepository.GetRoomByID(chatService.db, roomID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: chatService.constants.Field.Room}
		panic(notFoundError)
	}
	chatService.validateRoomParticipantAccess(senderID, room.CustomerID, room.CorporationID)
	message := &entity.ChatMessage{
		RoomID:   roomID,
		SenderID: senderID,
		Content:  content,
	}
	if err := chatService.chatRepository.CreateMessage(chatService.db, message); err != nil {
		panic(err)
	}
}

func (chatService *ChatService) GetRoomMessages(request chatdto.GetRoomMessageRequest) []chatdto.RoomMessagesResponse {
	chatService.userService.GetUserCredential(request.UserID)
	room, exist := chatService.chatRepository.GetRoomByID(chatService.db, request.RoomID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: chatService.constants.Field.Room}
		panic(notFoundError)
	}
	chatService.validateRoomParticipantAccess(request.UserID, room.CustomerID, room.CorporationID)
	messages := chatService.chatRepository.GetRoomMessages(chatService.db, request.RoomID)
	messagesResponse := make([]chatdto.RoomMessagesResponse, len(messages))
	for i, message := range messages {
		sender := chatService.userService.GetUserCredential(message.SenderID)
		messagesResponse[i] = chatdto.RoomMessagesResponse{
			Sender:  sender,
			Content: message.Content,
		}
	}
	return messagesResponse
}
