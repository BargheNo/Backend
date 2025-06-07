package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type ChatRepository struct{}

func NewChatRepository() *ChatRepository {
	return &ChatRepository{}
}

func (repo *ChatRepository) GetRoomByID(db database.Database, roomID uint) (*entity.ChatRoom, error) {
	var room *entity.ChatRoom
	result := db.GetDB().First(&room, roomID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return room, nil
}

func (repo *ChatRepository) GetUserRooms(db database.Database, userID uint) ([]*entity.ChatRoom, error) {
	var rooms []*entity.ChatRoom
	result := db.GetDB().Where("customer_id = ?", userID).Find(&rooms)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return rooms, nil
}

func (repo *ChatRepository) GetCorporationRooms(db database.Database, corporationID uint) ([]*entity.ChatRoom, error) {
	var rooms []*entity.ChatRoom
	result := db.GetDB().Where("corporation_id = ?", corporationID).Find(&rooms)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return rooms, nil
}

func (repo *ChatRepository) GetUserAndCorpRoom(db database.Database, userID uint, corporationID uint) (*entity.ChatRoom, error) {
	var room entity.ChatRoom
	result := db.GetDB().Where("customer_id = ? AND corporation_id = ?", userID, corporationID).First(&room)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &room, nil
}

func (repo *ChatRepository) GetRoomMessages(db database.Database, roomID uint, opts ...repository.QueryModifier) ([]*entity.ChatMessage, error) {
	var messages []*entity.ChatMessage
	query := db.GetDB().Where("room_id = ?", roomID)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

func (repo *ChatRepository) CreateRoom(db database.Database, room *entity.ChatRoom) error {
	return db.GetDB().Create(&room).Error
}

func (repo *ChatRepository) UpdateRoom(db database.Database, room *entity.ChatRoom) error {
	return db.GetDB().Save(&room).Error
}

func (repo *ChatRepository) CreateMessage(db database.Database, message *entity.ChatMessage) error {
	return db.GetDB().Create(&message).Error
}
