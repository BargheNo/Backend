package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type TicketComment struct {
	database.Model
	TicketID uint   `json:"ticket_id" gorm:"not null"`
	Ticket   Ticket `json:"ticket" gorm:"foreignKey:TicketID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OwnerID  uint   `json:"user_id" gorm:"not null"`
	Owner    User   `json:"user" gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Body     string `json:"body" gorm:"not null"`
}
