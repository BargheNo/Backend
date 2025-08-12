package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
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

func (ticketRepo *TicketRepository) GetCustomerTickets(db database.Database, ownerID uint, options *postgres.QueryOptions) ([]*entity.Ticket, error) {
	var tickets []*entity.Ticket
	query := db.GetDB().Where("owner_id = ?", ownerID)

	query = applyQueryOptions(query, options)
	result := query.Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

func (ticketRepo *TicketRepository) FindCustomerTicketsByQuery(db database.Database, ownerID uint, allowedStatus []enum.TicketStatus, query string, options *postgres.QueryOptions) ([]*entity.Ticket, error) {
	var tickets []*entity.Ticket
	result := db.GetDB().
		Where("owner_id = ? AND status IN ?", ownerID, allowedStatus).
		Where("description ILIKE ?", "%"+query+"%")

	result = applyQueryOptions(result, options)
	result = result.Find(&tickets)

	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

func (ticketRepo *TicketRepository) CountCustomerTicketsByQuery(db database.Database, ownerID uint, allowedStatus []enum.TicketStatus, query string) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Ticket{}).
		Where("owner_id = ? AND status IN ?", ownerID, allowedStatus).
		Where("description ILIKE ?", "%"+query+"%").
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ticketRepo *TicketRepository) GetTicketComments(db database.Database, ticketID uint, options *postgres.QueryOptions) ([]*entity.TicketComment, error) {
	var comments []*entity.TicketComment
	query := db.GetDB().Where("ticket_id = ?", ticketID)

	query = applyQueryOptions(query, options)
	result := query.Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

func (ticketRepo *TicketRepository) FindTicketByID(db database.Database, ticketID uint) (*entity.Ticket, error) {
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

func (ticketRepo *TicketRepository) GetTickets(db database.Database, options *postgres.QueryOptions) ([]*entity.Ticket, error) {
	var tickets []*entity.Ticket
	query := db.GetDB()

	query = applyQueryOptions(query, options)

	result := query.Find(&tickets)

	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

func (ticketRepo TicketRepository) FindTicketsByStatus(db database.Database, statuses []enum.TicketStatus, options *postgres.QueryOptions) ([]*entity.Ticket, error) {
	var tickets []*entity.Ticket
	query := db.GetDB().Where("status IN (?)", statuses)

	query = applyQueryOptions(query, options)

	result := query.Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}

	return tickets, nil
}

func (ticketRepo TicketRepository) CountTicketsByStatus(db database.Database, statuses []enum.TicketStatus) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Ticket{}).
		Where("status IN (?)", statuses).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (ticketRepo TicketRepository) FindTicketsByQuery(db database.Database, query string, options *postgres.QueryOptions) ([]*entity.Ticket, error) {
	var tickets []*entity.Ticket
	result := db.GetDB().
		Where("description ILIKE ?", "%"+query+"%")
	result = applyQueryOptions(result, options)
	result = result.Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	return tickets, nil
}

func (ticketRepo TicketRepository) CountTicketsByQuery(db database.Database, query string) (int64, error) {
	var count int64
	err := db.GetDB().
		Model(&entity.Ticket{}).
		Where("description ILIKE ?", "%"+query+"%").
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}
func (ticketRepo TicketRepository) FindCustomerTicketsByStatus(db database.Database, ownerID uint, statuses []enum.TicketStatus, options *postgres.QueryOptions) ([]*entity.Ticket, error) {
	var tickets []*entity.Ticket
	query := db.GetDB().Where("owner_id = ? AND status IN (?)", ownerID, statuses)

	query = applyQueryOptions(query, options)

	result := query.Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}

	return tickets, nil
}

func (ticketRepo TicketRepository) CountCustomerTicketsByStatus(db database.Database, ownerID uint, statuses []enum.TicketStatus) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Ticket{}).
		Where("owner_id = ? AND status IN (?)", ownerID, statuses).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}
