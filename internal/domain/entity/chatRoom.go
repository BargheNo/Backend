package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type ChatRoom struct {
	database.Model
	Type          enum.RoomType
	CorporationID uint
	CustomerID    uint
	Participants  []ChatParticipant
}
