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

func (chatService *ChatService) CreateChatRoom(request chatdto.CreateOrGetUserRoomRequest) (*entity.ChatRoom, error) {
	room := &entity.ChatRoom{
		CorporationID: request.CorporationID,
		CustomerID:    request.UserID,
		Status:        enum.ChatStatusActive,
	}
	err := chatService.chatRepository.CreateRoom(chatService.db, room)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (chatService *ChatService) CreateOrGetRoom(request chatdto.CreateOrGetUserRoomRequest) (chatdto.ChatRoomDetailsResponse, error) {
	customer, err := chatService.userService.GetUserCredential(request.UserID)
	if err != nil {
		return chatdto.ChatRoomDetailsResponse{}, err
	}
	corporation, err := chatService.corporationService.GetCorporationCredentials(request.CorporationID)
	if err != nil {
		return chatdto.ChatRoomDetailsResponse{}, err
	}
	var room *entity.ChatRoom
	room, err = chatService.chatRepository.GetUserAndCorpRoom(chatService.db, request.UserID, request.CorporationID)
	if err != nil {
		return chatdto.ChatRoomDetailsResponse{}, err
	}
	if room == nil {
		room, err = chatService.CreateChatRoom(request)
		if err != nil {
			return chatdto.ChatRoomDetailsResponse{}, err
		}
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

	return roomDetails, nil
}

func (chatService *ChatService) GetCorporationRoom(request chatdto.GetCorporationRoomRequest) (chatdto.ChatRoomDetailsResponse, error) {
	customerModel, err := chatService.userService.FindActiveUserByPhone(request.UserPhone)
	if err != nil {
		return chatdto.ChatRoomDetailsResponse{}, err
	}
	customerCred, err := chatService.userService.GetUserCredential(customerModel.ID)
	if err != nil {
		return chatdto.ChatRoomDetailsResponse{}, err
	}
	corporation, err := chatService.corporationService.GetCorporationCredentials(request.CorporationID)
	if err != nil {
		return chatdto.ChatRoomDetailsResponse{}, err
	}
	var room *entity.ChatRoom
	room, err = chatService.chatRepository.GetUserAndCorpRoom(chatService.db, customerModel.ID, request.CorporationID)
	if err != nil {
		return chatdto.ChatRoomDetailsResponse{}, err
	}
	if room == nil {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: chatService.constants.Field.Room,
		}
		return chatdto.ChatRoomDetailsResponse{}, forbiddenError
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

	return roomDetails, nil
}

func (chatService *ChatService) GetUserRooms(userID uint) ([]chatdto.ChatRoomDetailsResponse, error) {
	customer, err := chatService.userService.GetUserCredential(userID)
	if err != nil {
		return nil, err
	}
	rooms, err := chatService.chatRepository.GetUserRooms(chatService.db, userID)
	if err != nil {
		return nil, err
	}
	roomsDetails := make([]chatdto.ChatRoomDetailsResponse, len(rooms))
	for i, room := range rooms {
		corporation, err := chatService.corporationService.GetCorporationCredentials(room.CorporationID)
		if err != nil {
			return nil, err
		}
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
	return roomsDetails, nil
}

func (chatService *ChatService) GetCorporationRooms(request chatdto.GetCorporationRoomsRequest) ([]chatdto.ChatRoomDetailsResponse, error) {
	corporation, err := chatService.corporationService.GetCorporationCredentials(request.CorporationID)
	if err != nil {
		return nil, err
	}
	chatService.userService.DoesUserExist(request.ApplicantID)
	err = chatService.corporationService.CheckApplicantAccess(request.CorporationID, request.ApplicantID)
	if err != nil {
		return nil, err
	}
	rooms, err := chatService.chatRepository.GetCorporationRooms(chatService.db, request.CorporationID)
	if err != nil {
		return nil, err
	}
	roomsDetails := make([]chatdto.ChatRoomDetailsResponse, len(rooms))
	for i, room := range rooms {
		customer, err := chatService.userService.GetUserCredential(room.CustomerID)
		if err != nil {
			return nil, err
		}
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
	return roomsDetails, nil
}

func (chatService *ChatService) validateRoomParticipantAccess(senderID, memberID, corporationID uint) {
	if senderID != memberID {
		err := chatService.corporationService.CheckApplicantAccess(corporationID, senderID)
		if err != nil {
			return
		}
	}
}

func (chatService *ChatService) SaveMessage(roomID, senderID uint, content string) (chatdto.RoomMessagesResponse, error) {
	exist, err := chatService.userService.IsUserActive(senderID)
	if err != nil {
		return chatdto.RoomMessagesResponse{}, err
	}
	if !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: chatService.constants.Field.Room,
		}
		return chatdto.RoomMessagesResponse{}, forbiddenError
	}

	room, err := chatService.chatRepository.GetRoomByID(chatService.db, roomID)
	if err != nil {
		return chatdto.RoomMessagesResponse{}, err
	}
	if room == nil {
		notFoundError := exception.NotFoundError{Item: chatService.constants.Field.Room}
		return chatdto.RoomMessagesResponse{}, notFoundError
	}
	if room.Status == enum.ChatStatusBlocked {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: chatService.constants.Field.Room,
		}
		return chatdto.RoomMessagesResponse{}, forbiddenError
	}
	chatService.validateRoomParticipantAccess(senderID, room.CustomerID, room.CorporationID)

	message := &entity.ChatMessage{
		RoomID:   roomID,
		SenderID: senderID,
		Content:  content,
	}
	if err := chatService.chatRepository.CreateMessage(chatService.db, message); err != nil {
		return chatdto.RoomMessagesResponse{}, err
	}

	sender, err := chatService.userService.GetUserCredential(message.SenderID)
	if err != nil {
		return chatdto.RoomMessagesResponse{}, err
	}
	return chatdto.RoomMessagesResponse{
		ID:        message.ID,
		Sender:    sender,
		Content:   message.Content,
		TimeStamp: message.CreatedAt,
	}, nil
}

func (chatService *ChatService) GetRoomMessages(request chatdto.GetRoomMessageRequest) ([]chatdto.RoomMessagesResponse, error) {
	chatService.userService.DoesUserExist(request.UserID)
	room, err := chatService.chatRepository.GetRoomByID(chatService.db, request.RoomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		notFoundError := exception.NotFoundError{Item: chatService.constants.Field.Room}
		return nil, notFoundError
	}
	chatService.validateRoomParticipantAccess(request.UserID, room.CustomerID, room.CorporationID)
	paginationModifier := repositoryimpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	messages, err := chatService.chatRepository.GetRoomMessages(chatService.db, request.RoomID, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	messagesResponse := make([]chatdto.RoomMessagesResponse, len(messages))
	for i, message := range messages {
		sender, err := chatService.userService.GetUserCredential(message.SenderID)
		if err != nil {
			return nil, err
		}
		messagesResponse[i] = chatdto.RoomMessagesResponse{
			ID:        message.ID,
			Sender:    sender,
			Content:   message.Content,
			TimeStamp: message.CreatedAt,
		}
	}
	return messagesResponse, nil
}

func (chatService *ChatService) BlockChatRoom(request chatdto.BlockServiceChatRequest) error {
	chatService.userService.DoesUserExist(request.UserID)
	room, err := chatService.chatRepository.GetRoomByID(chatService.db, request.RoomID)
	if err != nil {
		return err
	}
	if room == nil {
		notFoundError := exception.NotFoundError{Item: chatService.constants.Field.Room}
		return notFoundError
	}
	if room.Status == enum.ChatStatusBlocked {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(chatService.constants.Field.Room, chatService.constants.Tag.AlreadyBlocked)
		return conflictErrors
	}
	chatService.validateRoomParticipantAccess(request.UserID, room.CustomerID, room.CorporationID)

	room.BlockedBy = &request.BlockedBy
	room.Status = request.ChatStatus
	err = chatService.chatRepository.UpdateRoom(chatService.db, room)
	if err != nil {
		return err
	}
	return nil
}

func (chatService *ChatService) UnBlockChatRoom(request chatdto.BlockServiceChatRequest) error {
	chatService.userService.DoesUserExist(request.UserID)
	room, err := chatService.chatRepository.GetRoomByID(chatService.db, request.RoomID)
	if err != nil {
		return err
	}
	if room == nil {
		notFoundError := exception.NotFoundError{Item: chatService.constants.Field.Room}
		return notFoundError
	}
	if room.Status == enum.ChatStatusActive {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(chatService.constants.Field.Room, chatService.constants.Tag.AlreadyActive)
		return conflictErrors
	}
	chatService.validateRoomParticipantAccess(request.UserID, room.CustomerID, room.CorporationID)

	if *room.BlockedBy != request.BlockedBy {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: chatService.constants.Field.Room,
		}
		return forbiddenError
	}
	room.BlockedBy = nil
	room.Status = request.ChatStatus
	err = chatService.chatRepository.UpdateRoom(chatService.db, room)
	if err != nil {
		return err
	}
	return nil
}
