package ticketdto

import (
	"mime/multipart"

	"github.com/BargheNo/Backend/internal/domain/enum"
)

type CreateTicketRequest struct {
	OwnerID     uint
	OwnerType   string
	Subject     enum.TicketSubject
	Description string
	Image       *multipart.FileHeader
}

type CreateCorporationTicketRequest struct {
	OperatorID    uint
	CorporationID uint
	Subject       enum.TicketSubject
	Description   string
	Image         *multipart.FileHeader
}

type TicketListRequest struct {
	OwnerID uint
	Status  uint
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
