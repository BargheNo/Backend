package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Ticket struct {
	database.Model
	OwnerID     uint               `gorm:"index"`
	Owner       User               `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Subject     enum.TicketSubject `gorm:"not null;index"`
	Description string             `gorm:"type:text;not null"`
	Image       string             `gorm:"type:varchar(255)"`
	Status      enum.TicketStatus  `gorm:"not null;index"`
}
