package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type TicketRepository struct {
}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{}
}

func (ticketRepo *TicketRepository) CreateTicket(db database.Database, ticket *entity.Ticket) error {
	return db.GetDB().Create(ticket).Error
}

func (ticketRepo *TicketRepository) UpdateTicket(db database.Database, ticket *entity.Ticket) error {
	return db.GetDB().Save(ticket).Error
}

func (ticketRepo *TicketRepository) GetCustomerTickets(db database.Database, ownerID uint) []*entity.Ticket {
	var tickets []*entity.Ticket
	db.GetDB().Where("owner_id = ?", ownerID).Find(&tickets)
	return tickets
}
