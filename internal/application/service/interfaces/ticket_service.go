package service

import ticketdto "github.com/BargheNo/Backend/internal/application/dto/ticket"

type TicketService interface {
	CreateTicket(requestInfo ticketdto.CreateTicketRequest)
	GetCustomerTickets(requestInfo ticketdto.TicketListRequest) []ticketdto.TicketResponse
	GetTicketComments(requestInfo ticketdto.TicketCommentListRequest) []ticketdto.TicketCommentResponse
	CreateTicketComment(requestInfo ticketdto.CreateTicketCommentRequest)
	GetTickets(requestInfo ticketdto.TicketListRequest) []ticketdto.TicketResponse
}
