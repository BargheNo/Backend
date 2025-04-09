package service

import (
	notificationdto "github.com/BargheNo/Backend/internal/application/dto/notification"
	"github.com/BargheNo/Backend/internal/domain/entity"
)

type NotificationService interface {
	CreateNotification(typeID, senderID, recipientID uint, additionalData map[string]interface{}) (uint, error)
	MarkAsRead(notificationdto.NotificationInfoRequest) error
	GetUserNotifications(userID uint, limit, offset int) []notificationdto.NotificationListResponse
	GetUserNotificationSettings(userID uint) []notificationdto.NotificationSettingResponse
	UpdateNotificationSettings(newSettingInfo notificationdto.UpdateSettingsRequest) error
	SendNotification(notification *entity.Notification) error
}
