package ticketdto

import (
	"mime/multipart"
)

type CreateTicketRequest struct {
	OwnerID     uint
	OwnerType   string
	Subject     uint
	Description string
	Image       *multipart.FileHeader
}

type TicketListRequest struct {
	OwnerID uint
	Status  uint
	Query   string
	Offset  int
	Limit   int
	SortBy  uint
	Asc     bool
}

type SearchTicketsRequest struct {
	Status uint
	Query  string
	Offset int
	Limit  int
	SortBy uint
	Asc    bool
}

type TicketCommentListRequest struct {
	TicketID uint
	OwnerID  uint
}

type CreateTicketCommentRequest struct {
	TicketID  uint
	OwnerID   uint
	OwnerType string
	Body      string
}

type ResolveTicketRequest struct {
	TicketID uint
	OwnerID  uint
}
