package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type NotificationRepository struct{}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

func (repo *NotificationRepository) GetNotificationByID(db database.Database, notificationID uint) (*entity.Notification, error) {
	var notification *entity.Notification
	result := db.GetDB().First(&notification, notificationID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return notification, nil
}

func (repo *NotificationRepository) GetNotificationsByTypesAndUserID(db database.Database, userID uint, types []uint, options *postgres.QueryOptions) ([]*entity.Notification, error) {
	var notifications []*entity.Notification
	query := db.GetDB().Where("recipient_id = ? and type_id IN ?", userID, types)
	query = applyQueryOptions(query, options)
	result := query.Find(&notifications)

	if result.Error != nil {
		return nil, result.Error
	}
	return notifications, nil
}

func (repo *NotificationRepository) CountNotificationsByTypesAndUserID(db database.Database, userID uint, types []uint) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Notification{}).
		Where("recipient_id = ? and type_id IN ?", userID, types).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *NotificationRepository) GetNotificationSettingByID(db database.Database, settingID uint) (*entity.NotificationSetting, error) {
	var setting *entity.NotificationSetting
	result := db.GetDB().First(&setting, settingID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return setting, nil
}

func (repo *NotificationRepository) GetNotificationSettingByUserAndType(db database.Database, userID, typeID uint) (*entity.NotificationSetting, error) {
	var setting *entity.NotificationSetting
	result := db.GetDB().Where("user_id = ? and type_id = ?", userID, typeID).First(&setting)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return setting, nil
}

func (repo *NotificationRepository) GetNotificationSettingByUserID(db database.Database, userID uint) ([]*entity.NotificationSetting, error) {
	var settings []*entity.NotificationSetting
	result := db.GetDB().Where("user_id = ?", userID).Find(&settings)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return settings, nil
}

func (repo *NotificationRepository) GetNotificationTypes(db database.Database) ([]*entity.NotificationType, error) {
	var notificationTypes []*entity.NotificationType
	result := db.GetDB().Find(&notificationTypes)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return notificationTypes, nil
}

func (repo *NotificationRepository) GetNotificationTypeByID(db database.Database, typeID uint) (*entity.NotificationType, error) {
	var notificationType *entity.NotificationType
	result := db.GetDB().First(&notificationType, typeID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return notificationType, nil
}

func (repo *NotificationRepository) GetNotificationTypeByName(db database.Database, name enum.NotificationType) (*entity.NotificationType, error) {
	var notificationType *entity.NotificationType
	result := db.GetDB().Where("name = ?", name).First(&notificationType)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return notificationType, nil
}

func (repo *NotificationRepository) CreateNotification(db database.Database, notification *entity.Notification) error {
	return db.GetDB().Create(&notification).Error
}

func (repo *NotificationRepository) UpdateNotification(db database.Database, notification *entity.Notification) error {
	return db.GetDB().Save(&notification).Error
}

func (repo *NotificationRepository) CreateNotificationSetting(db database.Database, setting *entity.NotificationSetting) error {
	return db.GetDB().Create(&setting).Error
}

func (repo *NotificationRepository) UpdateNotificationSetting(db database.Database, setting *entity.NotificationSetting) error {
	return db.GetDB().Save(&setting).Error
}

func (repo *NotificationRepository) CreateNotificationType(db database.Database, notificationType *entity.NotificationType) error {
	return db.GetDB().Create(&notificationType).Error
}
