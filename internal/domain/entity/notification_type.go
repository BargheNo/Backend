package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type NotificationType struct {
	database.Model
	Name        string `gorm:"type:varchar(100);not null"` // notif type will be enum
	Description string `gorm:"type:text"`
}
