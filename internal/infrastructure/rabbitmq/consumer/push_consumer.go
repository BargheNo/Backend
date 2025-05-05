package consumer

import (
	"encoding/json"
	"log"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/infrastructure/rabbitmq"
	"github.com/BargheNo/Backend/internal/infrastructure/websocket"
)

type PushConsumer struct {
	constants *bootstrap.RabbitMQConstants
	rabbitMQ  *rabbitmq.RabbitMQ
	wsHub     *websocket.Hub
}

func NewPushConsumer(
	constants *bootstrap.RabbitMQConstants,
	rabbitMQ *rabbitmq.RabbitMQ,
	wsHub *websocket.Hub,
) *PushConsumer {
	return &PushConsumer{
		constants: constants,
		rabbitMQ:  rabbitMQ,
		wsHub:     wsHub,
	}
}

func (consumer *PushConsumer) Start() error {
	return consumer.rabbitMQ.ConsumeMessages(consumer.constants.Events.NotificationsPush, consumer.handleMessage)
}

func (consumer *PushConsumer) handleMessage(body []byte) error {
	var msg struct {
		ID             uint      `json:"id"`
		Timestamp      time.Time `json:"timestamp"`
		Description    string    `json:"description"`
		Type           string    `json:"type"`
		AdditionalData string    `json:"additionalData"`
		IsRead         bool      `json:"isRead"`
		RecipientID    uint      `json:"recipientID"`
	}
	if err := json.Unmarshal(body, &msg); err != nil {
		log.Printf("Failed to unmarshal push notification message: %v", err)
	}

	payload := websocket.NotificationPayload{
		ID:             msg.ID,
		CreatedAt:      msg.Timestamp.Format(time.RFC3339),
		Type:           msg.Type,
		Description:    msg.Description,
		AdditionalData: msg,
		IsRead:         msg.IsRead,
	}

	content, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	consumer.wsHub.SendToUser(msg.RecipientID, websocket.MessageTypeNotification, content)

	return nil
}
