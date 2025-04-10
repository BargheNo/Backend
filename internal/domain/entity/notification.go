package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type Notification struct {
	database.Model
	TypeID         uint             `gorm:"not null;index"`
	Type           NotificationType `gorm:"foreignKey:TypeID"`
	AdditionalData string           `gorm:"type:text"`
	IsRead         bool             `gorm:"default:false"`
	RecipientID    uint             `gorm:"not null;index"`
}
