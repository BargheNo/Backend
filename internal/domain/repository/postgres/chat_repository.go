package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type ChatRepository interface {
	GetRoomByID(db database.Database, roomID uint) (*entity.ChatRoom, error)
	GetUserRooms(db database.Database, userID uint) ([]*entity.ChatRoom, error)
	GetCorporationRooms(db database.Database, corporationID uint) ([]*entity.ChatRoom, error)
	GetUserAndCorpRoom(db database.Database, userID uint, corporationID uint) (*entity.ChatRoom, error)
	GetRoomMessages(db database.Database, roomID uint, opts ...QueryModifier) ([]*entity.ChatMessage, error)
	CreateRoom(db database.Database, room *entity.ChatRoom) error
	UpdateRoom(db database.Database, room *entity.ChatRoom) error
	CreateMessage(db database.Database, message *entity.ChatMessage) error
}
