package service

import (
	notificationdto "github.com/BargheNo/Backend/internal/application/dto/notification"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
)

type NotificationService interface {
	CreateAndSendNotification(typeName enum.NotificationType, recipientID uint, additionalData interface{}) error
	CreateNotificationSettings(userID uint)
	GetUserNotificationSettings(userID uint) []notificationdto.NotificationSettingResponse
	GetUserNotifications(userID uint) []notificationdto.NotificationListResponse
	MarkAsRead(notificationInfo notificationdto.NotificationInfoRequest)
	SendNotification(notification *entity.Notification, notificationType *entity.NotificationType) error
	UpdateNotificationSettings(newSettingInfo notificationdto.UpdateSettingsRequest)
}
