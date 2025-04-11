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

var notificationTypes = map[enum.NotificationType]string{
	enum.ChatNotificationType: "شما یک پیام جدید دارید.",
}

func (seeder *NotificationTypeSeeder) SeedNotificationTypes() {
	for name, description := range notificationTypes {
		_, exist := seeder.notificationRepository.GetNotificationTypeByName(seeder.db, name)
		if !exist {
			notificationType := &entity.NotificationType{
				Name:        name,
				Description: description,
			}
			err := seeder.notificationRepository.CreateNotificationType(seeder.db, notificationType)
			if err != nil {
				panic(err)
			}
		}

	}
}
