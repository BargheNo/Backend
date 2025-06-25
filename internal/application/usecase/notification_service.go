package usecase

import (
	notificationdto "github.com/BargheNo/Backend/internal/application/dto/notification"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
)

type NotificationService interface {
	CreateAndSendNotification(typeName enum.NotificationType, recipientID uint, data []byte) error
	CreateNotificationSettings(userID uint) error
	GetNotificationsType() ([]notificationdto.NotificationTypeResponse, error)
	GetUserNotificationSettings(userID uint) ([]notificationdto.NotificationSettingResponse, error)
	GetUserNotifications(notificationsRequest notificationdto.NotificationListRequest) ([]notificationdto.NotificationListResponse, error)
	MarkAsRead(notificationInfo notificationdto.NotificationInfoRequest) error
	SendNotification(notification *entity.Notification, notificationType *entity.NotificationType) error
	UpdateNotificationSettings(newSettingInfo notificationdto.UpdateSettingsRequest) error
}
