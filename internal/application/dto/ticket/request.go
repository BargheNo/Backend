package ticketdto

import (
	"mime/multipart"

	"github.com/BargheNo/Backend/internal/domain/enum"
)

type CreateTicketRequest struct {
	OwnerID     uint
	Subject     enum.TicketSubject
	Description string
	Image       *multipart.FileHeader
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
	TicketID uint
	OwnerID  uint
	Body     string
	IsAdmin  bool
}

type ResolveTicketRequest struct {
	TicketID uint
	OwnerID  uint
}
