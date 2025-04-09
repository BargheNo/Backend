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

type ChatPayload struct {
	MessageID uint   `json:"message_id,omitempty"`
	Content   string `json:"content"`
}

type NotificationPayload struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	IsRead      bool   `json:"is_read"`
	CreatedAt   string `json:"created_at"`
}

// type Message struct {
// 	RoomID  uint
// 	Content []byte
// 	Client  *Client
// }
