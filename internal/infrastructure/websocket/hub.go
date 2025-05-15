package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type Hub struct {
	clients    map[uint]map[*Client]bool
	rooms      map[uint]map[*Client]bool
	broadcast  chan *Message
	Register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]map[*Client]bool),
		rooms:      make(map[uint]map[*Client]bool),
		broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			hub.handleRegister(client)
		case client := <-hub.unregister:
			hub.handleUnregister(client)
		case message := <-hub.broadcast:
			hub.handleBroadcast(message)
		}
	}
}

func (hub *Hub) handleRegister(client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if _, ok := hub.clients[client.userID]; !ok {
		hub.clients[client.userID] = make(map[*Client]bool)
	}
	hub.clients[client.userID][client] = true

	if _, ok := hub.rooms[client.roomID]; !ok {
		hub.rooms[client.roomID] = make(map[*Client]bool)
	}
	hub.rooms[client.roomID][client] = true
}

func (hub *Hub) handleUnregister(client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if clients, ok := hub.clients[client.userID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(hub.clients, client.userID)
		}
	}

	if room, ok := hub.rooms[client.roomID]; ok {
		delete(room, client)
		if len(room) == 0 {
			delete(hub.rooms, client.roomID)
		}
	}

	client.CloseConnection()
}

func (hub *Hub) handleBroadcast(message *Message) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	switch message.Type {
	case MessageTypeChat:
		messageBytes, err := json.Marshal(message)
		if err != nil {
			panic(err)
		}
		if room, ok := hub.rooms[message.RoomID]; ok {
			for client := range room {
				select {
				case client.send <- messageBytes:
				default:
					hub.unregister <- client
				}
			}
		}
	case MessageTypeNotification:
		if clients, ok := hub.clients[message.SenderID]; ok {
			for client := range clients {
				select {
				case client.send <- message.Content:
				default:
					hub.unregister <- client
				}
			}
		}
	}
}

func (hub *Hub) SendToUser(userID uint, messageType string, content []byte) error {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	clients, ok := hub.clients[userID]
	if !ok || len(clients) == 0 {
		return fmt.Errorf("no active connections for user %d", userID)
	}

	var safeClients []<-chan struct{}
	for client := range clients {
		select {
		case <-client.done:
			hub.mu.RUnlock()
			hub.mu.Lock()
			delete(clients, client)
			if len(clients) == 0 {
				delete(hub.clients, userID)
			}
			hub.mu.Unlock()
			hub.mu.RLock()
		default:
			safeClients = append(safeClients, client.done)
		}
	}

	if len(safeClients) == 0 {
		return fmt.Errorf("all connections for user %d are closed", userID)
	}

	for client := range clients {
		select {
		case <-client.done:
			continue
		case client.send <- content:
		default:
			select {
			case hub.unregister <- client:
			default:
				log.Printf("Failed to unregister client: unregister channel is full")
			}
		}
	}
	return nil
}
