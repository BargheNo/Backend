package consumer

import (
	"encoding/json"
	"log"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/application/port"
	"github.com/BargheNo/Backend/internal/infrastructure/rabbitmq"
)

type RegisterConsumer struct {
	constants           *bootstrap.RabbitMQConstants
	rabbitMQ            *rabbitmq.RabbitMQ
	notificationService port.NotificationService
}

func NewRegisterConsumer(
	constants *bootstrap.RabbitMQConstants,
	rabbitMQ *rabbitmq.RabbitMQ,
	notificationService port.NotificationService,
) *RegisterConsumer {
	return &RegisterConsumer{
		constants:           constants,
		rabbitMQ:            rabbitMQ,
		notificationService: notificationService,
	}
}

func (consumer *RegisterConsumer) Start() error {
	return consumer.rabbitMQ.ConsumeMessages(consumer.constants.Events.UserRegistered, consumer.handleMessage)
}

func (consumer *RegisterConsumer) handleMessage(body []byte) error {
	var msg struct {
		UserID uint `json:"userID"`
	}
	if err := json.Unmarshal(body, &msg); err != nil {
		log.Printf("Failed to unmarshal push notification message: %v", err)
	}

	consumer.notificationService.CreateNotificationSettings(msg.UserID)

	return nil
}
