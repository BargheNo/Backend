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
	Offset  int
	Limit   int
}

type TicketCommentListRequest struct {
	TicketID uint
	OwnerID  uint
	Offset   int
	Limit    int
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
