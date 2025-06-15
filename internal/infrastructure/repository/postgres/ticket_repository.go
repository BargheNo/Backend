package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
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

func (ticketRepo *TicketRepository) GetCustomerTickets(db database.Database, ownerID uint, opts ...repository.QueryModifier) ([]*entity.Ticket, error) {
	var tickets []*entity.Ticket
	query := db.GetDB().Where("owner_id = ?", ownerID)

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

func (ticketRepo *TicketRepository) GetTicketComments(db database.Database, ticketID uint, opts ...repository.QueryModifier) ([]*entity.TicketComment, error) {
	var comments []*entity.TicketComment
	query := db.GetDB().Where("ticket_id = ?", ticketID)

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

func (ticketRepo *TicketRepository) GetTicketByID(db database.Database, ticketID uint) (*entity.Ticket, error) {
	var ticket entity.Ticket
	result := db.GetDB().Where("id = ?", ticketID).First(&ticket)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &ticket, nil
}

func (ticketRepo *TicketRepository) CreateTicketComment(db database.Database, comment *entity.TicketComment) error {
	return db.GetDB().Create(comment).Error
}

func (ticketRepo *TicketRepository) GetTickets(db database.Database, opts ...repository.QueryModifier) ([]*entity.Ticket, error) {
	var tickets []*entity.Ticket
	query := db.GetDB()

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}

	result := query.Find(&tickets)

	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}
