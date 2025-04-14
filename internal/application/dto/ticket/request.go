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
