package service

import ticketdto "github.com/BargheNo/Backend/internal/application/dto/ticket"

type TicketService interface {
	CreateCustomerTicket(requestInfo ticketdto.CreateTicketRequest)
	GetCustomerTickets(requestInfo ticketdto.TicketListRequest) []ticketdto.TicketResponse
	CreateCustomerTicketComment(requestInfo ticketdto.CreateTicketCommentRequest)
	GetCustomerTicketComments(requestInfo ticketdto.TicketCommentListRequest) []ticketdto.TicketCommentResponse
	CreateAdminTicketComment(requestInfo ticketdto.CreateTicketCommentRequest)
	GetAdminTickets(requestInfo ticketdto.TicketListRequest) []ticketdto.TicketResponse
	GetAdminTicketComments(requestInfo ticketdto.TicketCommentListRequest) []ticketdto.TicketCommentResponse
	ResolveTicket(requestInfo ticketdto.ResolveTicketRequest)
}
