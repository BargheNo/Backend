package consumer

import (
	"encoding/json"
	"log"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/rabbitmq"
)

type SendNotificationConsumer struct {
	constants           *bootstrap.RabbitMQConstants
	rabbitMQ            *rabbitmq.RabbitMQ
	notificationService usecase.NotificationService
}

func NewSendNotificationConsumer(
	constants *bootstrap.RabbitMQConstants,
	rabbitMQ *rabbitmq.RabbitMQ,
	notificationService usecase.NotificationService,
) *SendNotificationConsumer {
	return &SendNotificationConsumer{
		constants:           constants,
		rabbitMQ:            rabbitMQ,
		notificationService: notificationService,
	}
}

func (consumer *SendNotificationConsumer) Start() error {
	return consumer.rabbitMQ.ConsumeMessages(consumer.constants.Events.SendNotification, consumer.handleMessage)
}

func (consumer *SendNotificationConsumer) handleMessage(body []byte) error {
	var msg struct {
		TypeName    enum.NotificationType `json:"typeName"`
		RecipientID uint                  `json:"recipientID"`
		Data        []byte                `json:"data"`
	}
	if err := json.Unmarshal(body, &msg); err != nil {
		log.Printf("Failed to unmarshal push notification message: %v", err)
	}

	if err := consumer.notificationService.CreateAndSendNotification(msg.TypeName, msg.RecipientID, msg.Data); err != nil {
		return err
	}

	return nil
}
