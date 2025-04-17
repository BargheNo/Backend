package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	chatdto "github.com/BargheNo/Backend/internal/application/dto/chat"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
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

func (chatService *ChatService) CreateChatRoom(request chatdto.CreateOrGetUserRoomRequest) *entity.ChatRoom {
	room := &entity.ChatRoom{
		CorporationID: request.CorporationID,
		CustomerID:    request.UserID,
		Status:        enum.ChatStatusActive,
	}
	err := chatService.chatRepository.CreateRoom(chatService.db, room)
	if err != nil {
		panic(err)
	}
	return room
}

func (chatService *ChatService) CreateOrGetRoom(request chatdto.CreateOrGetUserRoomRequest) chatdto.ChatRoomDetailsResponse {
	customer := chatService.userService.GetUserCredential(request.UserID)
	corporation := chatService.corporationService.GetCorporationCredentials(request.CorporationID)
	var room *entity.ChatRoom
	var exist bool
	room, exist = chatService.chatRepository.GetUserAndCorpRoom(chatService.db, request.UserID, request.CorporationID)
	if !exist {
		room = chatService.CreateChatRoom(request)
	}

	blockedBy := ""
	if room.Status == enum.ChatStatusBlocked {
		blockedBy = room.BlockedBy.String()
	}
	roomDetails := chatdto.ChatRoomDetailsResponse{
		RoomID:                room.ID,
		CustomerCredential:    customer,
		CorporationCredential: corporation,
		Status:                room.Status.String(),
		BlockedBy:             blockedBy,
	}

	return roomDetails
}

func (chatService *ChatService) GetCorporationRoom(request chatdto.GetCorporationRoomRequest) chatdto.ChatRoomDetailsResponse {
	customerModel := chatService.userService.FindUserByPhone(request.UserPhone)
	customerCred := chatService.userService.GetUserCredential(customerModel.ID)
	corporation := chatService.corporationService.GetCorporationCredentials(request.CorporationID)
	var room *entity.ChatRoom
	var exist bool
	room, exist = chatService.chatRepository.GetUserAndCorpRoom(chatService.db, customerModel.ID, request.CorporationID)
	if !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: chatService.constants.Field.Room,
		}
		panic(forbiddenError)
	}
	blockedBy := ""
	if room.Status == enum.ChatStatusBlocked {
		blockedBy = room.BlockedBy.String()
	}
	roomDetails := chatdto.ChatRoomDetailsResponse{
		RoomID:                room.ID,
		CustomerCredential:    customerCred,
		CorporationCredential: corporation,
		Status:                room.Status.String(),
		BlockedBy:             blockedBy,
	}

	return roomDetails
}

func (chatService *ChatService) GetUserRooms(userID uint) []chatdto.ChatRoomDetailsResponse {
	customer := chatService.userService.GetUserCredential(userID)
	rooms := chatService.chatRepository.GetUserRooms(chatService.db, userID)
	roomsDetails := make([]chatdto.ChatRoomDetailsResponse, len(rooms))
	for i, room := range rooms {
		corporation := chatService.corporationService.GetCorporationCredentials(room.CorporationID)
		blockedBy := ""
		if room.Status == enum.ChatStatusBlocked {
			blockedBy = room.BlockedBy.String()
		}
		roomsDetails[i] = chatdto.ChatRoomDetailsResponse{
			RoomID:                room.ID,
			CustomerCredential:    customer,
			CorporationCredential: corporation,
			Status:                room.Status.String(),
			BlockedBy:             blockedBy,
		}
	}
	return roomsDetails
}

func (chatService *ChatService) GetCorporationRooms(request chatdto.GetCorporationRoomsRequest) []chatdto.ChatRoomDetailsResponse {
	corporation := chatService.corporationService.GetCorporationCredentials(request.CorporationID)
	chatService.userService.DoesUserExist(request.ApplicantID)
	chatService.corporationService.CheckApplicantAccess(request.CorporationID, request.ApplicantID)
	rooms := chatService.chatRepository.GetCorporationRooms(chatService.db, request.CorporationID)
	roomsDetails := make([]chatdto.ChatRoomDetailsResponse, len(rooms))
	for i, room := range rooms {
		customer := chatService.userService.GetUserCredential(room.CustomerID)
		blockedBy := ""
		if room.Status == enum.ChatStatusBlocked {
			blockedBy = room.BlockedBy.String()
		}
		roomsDetails[i] = chatdto.ChatRoomDetailsResponse{
			RoomID:                room.ID,
			CustomerCredential:    customer,
			CorporationCredential: corporation,
			Status:                room.Status.String(),
			BlockedBy:             blockedBy,
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
	exist := chatService.userService.IsUserActive(senderID)
	if !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: chatService.constants.Field.Room,
		}
		panic(forbiddenError)
	}
	room, exist := chatService.chatRepository.GetRoomByID(chatService.db, roomID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: chatService.constants.Field.Room}
		panic(notFoundError)
	}
	if room.Status == enum.ChatStatusBlocked {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: chatService.constants.Field.Room,
		}
		panic(forbiddenError)
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
	chatService.userService.DoesUserExist(request.UserID)
	room, exist := chatService.chatRepository.GetRoomByID(chatService.db, request.RoomID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: chatService.constants.Field.Room}
		panic(notFoundError)
	}
	chatService.validateRoomParticipantAccess(request.UserID, room.CustomerID, room.CorporationID)
	paginationModifier := repositoryimpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	messages := chatService.chatRepository.GetRoomMessages(chatService.db, request.RoomID, paginationModifier, sortingModifier)
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

func (chatService *ChatService) BlockChatRoom(request chatdto.BlockServiceChatRequest) {
	chatService.userService.DoesUserExist(request.UserID)
	room, exist := chatService.chatRepository.GetRoomByID(chatService.db, request.RoomID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: chatService.constants.Field.Room}
		panic(notFoundError)
	}
	if room.Status == enum.ChatStatusBlocked {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(chatService.constants.Field.Room, chatService.constants.Tag.AlreadyBlocked)
		panic(conflictErrors)
	}
	chatService.validateRoomParticipantAccess(request.UserID, room.CustomerID, room.CorporationID)

	room.BlockedBy = &request.BlockedBy
	room.Status = request.ChatStatus
	err := chatService.chatRepository.UpdateRoom(chatService.db, room)
	if err != nil {
		panic(err)
	}
}

func (chatService *ChatService) UnBlockChatRoom(request chatdto.BlockServiceChatRequest) {
	chatService.userService.DoesUserExist(request.UserID)
	room, exist := chatService.chatRepository.GetRoomByID(chatService.db, request.RoomID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: chatService.constants.Field.Room}
		panic(notFoundError)
	}
	if room.Status == enum.ChatStatusActive {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(chatService.constants.Field.Room, chatService.constants.Tag.AlreadyActive)
		panic(conflictErrors)
	}
	chatService.validateRoomParticipantAccess(request.UserID, room.CustomerID, room.CorporationID)

	if *room.BlockedBy != request.BlockedBy {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: chatService.constants.Field.Room,
		}
		panic(forbiddenError)
	}
	room.BlockedBy = nil
	room.Status = request.ChatStatus
	err := chatService.chatRepository.UpdateRoom(chatService.db, room)
	if err != nil {
		panic(err)
	}
}
