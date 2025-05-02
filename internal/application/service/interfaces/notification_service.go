package service

import (
	notificationdto "github.com/BargheNo/Backend/internal/application/dto/notification"
	"github.com/BargheNo/Backend/internal/domain/entity"
)

type NotificationService interface {
	CreateAndSendNotification(typeID uint, recipientID uint, additionalData interface{}) error
	CreateNotificationSettings(userID uint)
	GetUserNotificationSettings(userID uint) []notificationdto.NotificationSettingResponse
	GetUserNotifications(userID uint) []notificationdto.NotificationListResponse
	MarkAsRead(notificationInfo notificationdto.NotificationInfoRequest)
	SendNotification(notification *entity.Notification) error
	UpdateNotificationSettings(newSettingInfo notificationdto.UpdateSettingsRequest)
}
