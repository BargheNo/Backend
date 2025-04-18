package websocket

import (
	"encoding/json"
	"time"
)

// convert to enum
const (
	MessageTypeChat         = "chat"
	MessageTypeNotification = "notification"
)

type Message struct {
	Type      string          `json:"type"`
	RoomID    uint            `json:"room_id,omitempty"`
	SenderID  uint            `json:"sender_id,omitempty"`
	Content   json.RawMessage `json:"content"`
	Timestamp time.Time       `json:"timestamp"`
	Client    *Client         `json:"-"`
}

type NotificationPayload struct {
	ID             uint              `json:"id"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	AdditionalData map[string]string `json:"additionalData"`
	Type           string            `json:"type"`
	IsRead         bool              `json:"is_read"`
	CreatedAt      string            `json:"created_at"`
}
