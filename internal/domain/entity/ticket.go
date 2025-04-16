package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Ticket struct {
	database.Model
	Subject     enum.TicketSubject `gorm:"not null;index"`
	Description string             `gorm:"type:text;not null"`
	Image       string             `gorm:"type:varchar(255)"`
	Status      enum.TicketStatus  `gorm:"not null;index"`
	OwnerID     uint               `gorm:"not null;index"`
	OwnerType   string             `gorm:"type:varchar(50);not null"`
}
