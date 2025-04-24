package websocket

import (
	"bytes"
	"encoding/json"
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
	closeOnce           sync.Once
	chatService         service.ChatService
	notificationService service.NotificationService
}

func NewClient(
	hub *Hub, conn any, roomID, userID uint,
	websocketSetting *bootstrap.WebsocketSetting,
	chatService service.ChatService,
	notificationService service.NotificationService,
) *Client {
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
	defer client.CloseConnection()

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
		client.CloseConnection()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.mu.Lock()
			client.conn.SetWriteDeadline(time.Now().Add(client.websocketSetting.WriteTimeout))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closed by server"))
				client.mu.Unlock()
				return nil
			}

			writer, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				client.mu.Unlock()
				return err
			}
			writer.Write(message)

			n := len(client.send)
			for i := 0; i < n; i++ {
				writer.Write(bytes.TrimSpace([]byte{'\n'}))
				writer.Write(<-client.send)
			}
			writer.Close()
			client.mu.Unlock()

		case <-ticker.C:
			client.mu.Lock()
			client.conn.SetWriteDeadline(time.Now().Add(client.websocketSetting.WriteTimeout))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				client.mu.Unlock()
				return err
			}
			client.mu.Unlock()

		case <-client.done:
			return nil
		}
	}
}

func (client *Client) CloseConnection() {
	client.closeOnce.Do(func() {
		close(client.done)
		close(client.send)
		client.conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing connection"))
		client.conn.Close()
	})
}

func (client *Client) processAndSaveChatMessage(message *Message) {
	var content string
	if err := json.Unmarshal(message.Content, &content); err != nil {
		return
	}
	savedMessage := client.chatService.SaveMessage(client.roomID, client.userID, content)
	message.MessageID = savedMessage.ID
	message.Sender = savedMessage.Sender
}
