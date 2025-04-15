package ticketdto

import (
	"time"

	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type TicketResponse struct {
	ID          uint `json:"id"`
	Owner       userdto.CredentialResponse
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
}

type TicketCommentResponse struct {
	ID     uint `json:"id"`
	Author TicketCommentAuthorResponse
	Body   string `json:"body"`
}

type TicketCommentAuthorResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	OwnerType string `json:"owner_type"`
}
