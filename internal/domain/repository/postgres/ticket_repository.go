package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type TicketRepository interface {
	CreateTicket(db database.Database, ticket *entity.Ticket) error
	GetCustomerTickets(db database.Database, ownerID uint, opts ...QueryModifier) ([]*entity.Ticket, error)
	UpdateTicket(db database.Database, ticket *entity.Ticket) error
	GetTicketComments(db database.Database, ticketID uint, opts ...QueryModifier) ([]*entity.TicketComment, error)
	FindTicketByID(db database.Database, ticketID uint) (*entity.Ticket, error)
	CreateTicketComment(db database.Database, comment *entity.TicketComment) error
	GetTickets(db database.Database, opts ...QueryModifier) ([]*entity.Ticket, error)
	FindTicketsByStatus(db database.Database, statuses []enum.TicketStatus, opts ...QueryModifier) ([]*entity.Ticket, error)
	FindCustomerTicketsByStatus(db database.Database, ownerID uint, statuses []enum.TicketStatus, opts ...QueryModifier) ([]*entity.Ticket, error)
}
