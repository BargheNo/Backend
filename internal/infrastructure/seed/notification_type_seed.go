package seed

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type NotificationTypeSeeder struct {
	notificationRepository repository.NotificationRepository
	db                     database.Database
}

func NewNotificationTypeSeeder(
	notificationRepository repository.NotificationRepository,
	db database.Database,
) *NotificationTypeSeeder {
	return &NotificationTypeSeeder{
		notificationRepository: notificationRepository,
		db:                     db,
	}
}

func (seeder *NotificationTypeSeeder) SeedNotificationTypes() {
	for _, notification := range enum.GetAllNotificationTypes() {
		_, exist := seeder.notificationRepository.GetNotificationTypeByName(seeder.db, notification)
		if !exist {
			notificationType := &entity.NotificationType{
				Name:              notification,
				Description:       notification.Description(),
				EmailTemplatePath: notification.EmailTemplatePath(),
			}
			err := seeder.notificationRepository.CreateNotificationType(seeder.db, notificationType)
			if err != nil {
				panic(err)
			}
		}

	}
}
