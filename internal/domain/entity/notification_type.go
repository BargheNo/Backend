package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type NotificationType struct {
	database.Model
	Name              enum.NotificationType `gorm:"unique;not null"`
	Description       string                `gorm:"type:text"`
	EmailTemplatePath string
}
