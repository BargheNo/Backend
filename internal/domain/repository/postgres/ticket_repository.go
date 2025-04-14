package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type TicketRepository interface {
	CreateTicket(db database.Database, ticket *entity.Ticket) error
	UpdateTicket(db database.Database, ticket *entity.Ticket) error
	GetCustomerTickets(db database.Database, ownerID uint) []*entity.Ticket
}
