package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type TicketRepository interface {
	CreateTicket(db database.Database, ticket *entity.Ticket) error
	GetCustomerTickets(db database.Database, ownerID uint, options *QueryOptions) ([]*entity.Ticket, error)
	FindCustomerTicketsByQuery(db database.Database, ownerID uint, allowedStatus []enum.TicketStatus, query string, options *QueryOptions) ([]*entity.Ticket, error)
	CountCustomerTicketsByQuery(db database.Database, ownerID uint, allowedStatus []enum.TicketStatus, query string) (int64, error)
	UpdateTicket(db database.Database, ticket *entity.Ticket) error
	GetTicketComments(db database.Database, ticketID uint, options *QueryOptions) ([]*entity.TicketComment, error)
	FindTicketByID(db database.Database, ticketID uint) (*entity.Ticket, error)
	CreateTicketComment(db database.Database, comment *entity.TicketComment) error
	GetTickets(db database.Database, options *QueryOptions) ([]*entity.Ticket, error)
	FindTicketsByStatus(db database.Database, statuses []enum.TicketStatus, subjects []enum.TicketSubject, options *QueryOptions) ([]*entity.Ticket, error)
	CountTicketsByStatus(db database.Database, statuses []enum.TicketStatus, subjects []enum.TicketSubject) (int64, error)
	FindTicketsByQuery(db database.Database, query string, allowedStatuses []enum.TicketStatus, allowedSubjects []enum.TicketSubject, options *QueryOptions) ([]*entity.Ticket, error)
	CountTicketsByQuery(db database.Database, query string, allowedStatuses []enum.TicketStatus, allowedSubjects []enum.TicketSubject) (int64, error)
	FindCustomerTicketsByStatus(db database.Database, ownerID uint, statuses []enum.TicketStatus, options *QueryOptions) ([]*entity.Ticket, error)
	CountCustomerTicketsByStatus(db database.Database, ownerID uint, statuses []enum.TicketStatus) (int64, error)
}
