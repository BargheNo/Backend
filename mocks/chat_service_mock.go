package mocks

import (
	chatdto "github.com/BargheNo/Backend/internal/application/dto/chat"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type ChatServiceMock struct {
	mock.Mock
}

func NewChatServiceMock() *ChatServiceMock {
	return &ChatServiceMock{}
}

func (m *ChatServiceMock) CreateChatRoom(request chatdto.CreateOrGetUserRoomRequest) *entity.ChatRoom {
	args := m.Called(request)
	if room, ok := args.Get(0).(*entity.ChatRoom); ok {
		return room
	}
	return nil
}

func (m *ChatServiceMock) CreateOrGetRoom(request chatdto.CreateOrGetUserRoomRequest) chatdto.ChatRoomDetailsResponse {
	args := m.Called(request)
	return args.Get(0).(chatdto.ChatRoomDetailsResponse)
}

func (m *ChatServiceMock) GetCorporationRoom(request chatdto.GetCorporationRoomRequest) chatdto.ChatRoomDetailsResponse {
	args := m.Called(request)
	return args.Get(0).(chatdto.ChatRoomDetailsResponse)
}

func (m *ChatServiceMock) GetUserRooms(userID uint) []chatdto.ChatRoomDetailsResponse {
	args := m.Called(userID)
	return args.Get(0).([]chatdto.ChatRoomDetailsResponse)
}

func (m *ChatServiceMock) GetCorporationRooms(request chatdto.GetCorporationRoomsRequest) []chatdto.ChatRoomDetailsResponse {
	args := m.Called(request)
	return args.Get(0).([]chatdto.ChatRoomDetailsResponse)
}

func (m *ChatServiceMock) SaveMessage(roomID, senderID uint, content string) chatdto.RoomMessagesResponse {
	args := m.Called(roomID, senderID, content)
	return args.Get(0).(chatdto.RoomMessagesResponse)
}

func (m *ChatServiceMock) GetRoomMessages(request chatdto.GetRoomMessageRequest) []chatdto.RoomMessagesResponse {
	args := m.Called(request)
	return args.Get(0).([]chatdto.RoomMessagesResponse)
}

func (m *ChatServiceMock) BlockChatRoom(request chatdto.BlockServiceChatRequest) {
	m.Called(request)
}

func (m *ChatServiceMock) UnBlockChatRoom(request chatdto.BlockServiceChatRequest) {
	m.Called(request)
}
