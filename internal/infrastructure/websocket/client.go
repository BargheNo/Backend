package websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/gorilla/websocket"
)

type Client struct {
	websocketSetting    *bootstrap.WebsocketSetting
	Hub                 *Hub
	conn                *websocket.Conn
	send                chan []byte
	roomID              uint
	userID              uint
	mu                  sync.Mutex
	done                chan struct{}
	chatService         service.ChatService
	notificationService service.NotificationService
}

func NewClient(
	hub *Hub, conn any, roomID, userID uint,
	websocketSetting *bootstrap.WebsocketSetting,
	chatService service.ChatService,
	notificationService service.NotificationService,
) *Client {
	if hub == nil {
		panic("hub cannot be nil")
	}
	wsConn, _ := conn.(*websocket.Conn)
	return &Client{
		websocketSetting:    websocketSetting,
		Hub:                 hub,
		conn:                wsConn,
		send:                make(chan []byte, websocketSetting.MessageBufferSize),
		roomID:              roomID,
		userID:              userID,
		done:                make(chan struct{}),
		chatService:         chatService,
		notificationService: notificationService,
	}
}

func (client *Client) ReadPump() error {
	defer func() {
		client.Hub.unregister <- client
		close(client.done)
		client.conn.Close()
	}()

	client.conn.SetReadLimit(int64(client.websocketSetting.MaxMessageSize))
	client.conn.SetReadDeadline(time.Now().Add(client.websocketSetting.ReadTimeout))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(client.websocketSetting.ReadTimeout))
		return nil
	})

	for {
		_, rawMessage, err := client.conn.ReadMessage()
		if err != nil {
			return err
		}

		var message Message
		if err := json.Unmarshal(rawMessage, &message); err != nil {
			continue
		}
		message.Client = client
		message.Timestamp = time.Now()
		message.RoomID = client.roomID

		switch message.Type {
		case MessageTypeChat:
			client.processAndSaveChatMessage(&message)
		}

		client.Hub.broadcast <- &message
	}
}

func (client *Client) WritePump() error {
	ticker := time.NewTicker(client.websocketSetting.PingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.mu.Lock()
			client.conn.SetWriteDeadline(time.Now().Add(client.websocketSetting.WriteTimeout))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				client.mu.Unlock()
				return fmt.Errorf("send channel closed")
			}

			writer, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				client.mu.Unlock()
				return fmt.Errorf("failed to get writer: %w", err)
			}
			writer.Write(message)

			n := len(client.send)
			for i := 0; i < n; i++ {
				writer.Write(bytes.TrimSpace([]byte{'\n'}))
				writer.Write(<-client.send)
			}

			if err := writer.Close(); err != nil {
				client.mu.Unlock()
				return fmt.Errorf("failed to close writer: %w", err)
			}
			client.mu.Unlock()

		case <-ticker.C:
			client.mu.Lock()
			client.conn.SetWriteDeadline(time.Now().Add(client.websocketSetting.WriteTimeout))
			err := client.conn.WriteMessage(websocket.PingMessage, nil)
			client.mu.Unlock()
			if err != nil {
				return err
			}

		case <-client.done:
			return fmt.Errorf("client connection done")
		}
	}
}

func (client *Client) processAndSaveChatMessage(message *Message) {
	var content string
	if err := json.Unmarshal(message.Content, &content); err != nil {
		panic(err)
	}
	client.chatService.SaveMessage(client.roomID, client.userID, content)
}

// func (client *Client) processAndSaveNotificationMessage(message *Message) error {
// }
