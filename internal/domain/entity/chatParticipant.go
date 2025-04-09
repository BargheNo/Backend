package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type ChatParticipant struct {
	database.Model
	RoomID          uint                 `gorm:"not null;index"`
	UserID          uint                 `gorm:"not null;index"`
	User            User                 `gorm:"foreignKey:UserID"`
	ParticipantType enum.ParticipantType `gorm:"not null"`
}
