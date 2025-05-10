package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type NotificationRepository interface {
	GetNotificationByID(db database.Database, notificationID uint) (*entity.Notification, bool)
	GetNotificationsByTypesAndUserID(db database.Database, userID uint, types []uint, opts ...QueryModifier) []*entity.Notification
	GetNotificationSettingByID(db database.Database, settingID uint) (*entity.NotificationSetting, bool)
	GetNotificationSettingByUserAndType(db database.Database, userID, typeID uint) (*entity.NotificationSetting, bool)
	GetNotificationSettingByUserID(db database.Database, userID uint) []*entity.NotificationSetting
	GetNotificationTypeByID(db database.Database, typeID uint) (*entity.NotificationType, bool)
	GetNotificationTypes(db database.Database) []*entity.NotificationType
	GetNotificationTypeByName(db database.Database, name enum.NotificationType) (*entity.NotificationType, bool)
	CreateNotification(db database.Database, notification *entity.Notification) error
	UpdateNotification(db database.Database, notification *entity.Notification) error
	CreateNotificationSetting(db database.Database, setting *entity.NotificationSetting) error
	UpdateNotificationSetting(db database.Database, setting *entity.NotificationSetting) error
	CreateNotificationType(db database.Database, notificationType *entity.NotificationType) error
}
