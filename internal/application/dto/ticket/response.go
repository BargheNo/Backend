package ticketdto

import (
	"time"

	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type TicketResponse struct {
	ID          uint                       `json:"id"`
	Owner       userdto.CredentialResponse `json:"owner"`
	Subject     string                     `json:"subject"`
	Description string                     `json:"description"`
	Status      string                     `json:"status"`
	Image       string                     `json:"image"`
	CreatedAt   time.Time                  `json:"created_at"`
}

type TicketCommentResponse struct {
	ID         uint                       `json:"id"`
	AuthorType string                     `json:"authorType"`
	Author     userdto.CredentialResponse `json:"author"`
	Body       string                     `json:"body"`
}

type TicketEnumResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
