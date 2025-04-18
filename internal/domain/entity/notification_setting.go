package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type NotificationSetting struct {
	database.Model
	UserID           uint             `gorm:"not null;index"`
	TypeID           uint             `gorm:"not null;index"`
	NotificationType NotificationType `gorm:"foreignKey:TypeID"`
	IsEmailEnabled   bool             `gorm:"default:true"`
	IsPushEnabled    bool             `gorm:"default:true"`
}
