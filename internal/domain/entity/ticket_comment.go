package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type TicketComment struct {
	database.Model
	OwnerID   uint   `gorm:"not null;index"`
	OwnerType string `gorm:"type:varchar(50);not null"`
	TicketID  uint   `gorm:"not null;index"`
	Body      string `gorm:"type:text;not null"`
}
