package usecase

import ticketdto "github.com/BargheNo/Backend/internal/application/dto/ticket"

type TicketService interface {
	GetTicketSortableColumns() []ticketdto.TicketEnumResponse
	CreateCustomerTicket(requestInfo ticketdto.CreateTicketRequest) error
	GetCustomerTickets(requestInfo ticketdto.TicketListRequest) ([]ticketdto.TicketResponse, int64, error)
	CreateCustomerTicketComment(requestInfo ticketdto.CreateTicketCommentRequest) error
	GetCustomerTicketComments(requestInfo ticketdto.TicketCommentListRequest) ([]ticketdto.TicketCommentResponse, error)
	CreateAdminTicketComment(requestInfo ticketdto.CreateTicketCommentRequest) error
	GetAdminTickets(requestInfo ticketdto.TicketListRequest) ([]ticketdto.TicketResponse, int64, error)
	SearchTickets(requestInfo ticketdto.TicketListRequest) ([]ticketdto.TicketResponse, int64, error)
	GetAdminTicketComments(requestInfo ticketdto.TicketCommentListRequest) ([]ticketdto.TicketCommentResponse, error)
	ResolveTicket(requestInfo ticketdto.ResolveTicketRequest) error
	GetTicketStatuses() []ticketdto.TicketEnumResponse
}
