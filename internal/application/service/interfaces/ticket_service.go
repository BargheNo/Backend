package service

import ticketdto "github.com/BargheNo/Backend/internal/application/dto/ticket"

type TicketService interface {
	CreateTicket(requestInfo ticketdto.CreateTicketRequest)
}
