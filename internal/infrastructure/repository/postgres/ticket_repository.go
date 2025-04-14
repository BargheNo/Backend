package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
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

func (ticketRepo *TicketRepository) GetTicketComments(db database.Database, ticketID uint) []*entity.TicketComment {
	var comments []*entity.TicketComment
	db.GetDB().Where("ticket_id = ?", ticketID).Find(&comments)
	return comments
}

func (ticketRepo *TicketRepository) GetTicketByID(db database.Database, ticketID uint) (*entity.Ticket, bool) {
	var ticket entity.Ticket
	result := db.GetDB().Where("id = ?", ticketID).First(&ticket)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &ticket, true
}

func (ticketRepo *TicketRepository) CreateTicketComment(db database.Database, comment *entity.TicketComment) error {
	return db.GetDB().Create(comment).Error
}
