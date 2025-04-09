package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type ChatMessage struct {
	database.Model
	RoomID   uint   `json:"room_id"`
	SenderID uint   `json:"sender_id"`
	Content  string `json:"content"`
}
