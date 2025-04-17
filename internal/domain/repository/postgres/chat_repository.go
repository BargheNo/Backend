package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type ChatRepository interface {
	GetRoomByID(db database.Database, roomID uint) (*entity.ChatRoom, bool)
	GetUserRooms(db database.Database, userID uint) []*entity.ChatRoom
	GetCorporationRooms(db database.Database, corporationID uint) []*entity.ChatRoom
	GetUserAndCorpRoom(db database.Database, userID uint, corporationID uint) (*entity.ChatRoom, bool)
	GetRoomMessages(db database.Database, roomID uint) []*entity.ChatMessage
	CreateRoom(db database.Database, room *entity.ChatRoom) error
	CreateMessage(db database.Database, message *entity.ChatMessage) error
}
