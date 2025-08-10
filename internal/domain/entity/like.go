package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type Like struct {
	database.Model
	UserID    uint   `gorm:"not null;index"`
	User      User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	OwnerID   uint   `gorm:"not null;index"`
	OwnerType string `gorm:"type:varchar(50);not null"`
}
