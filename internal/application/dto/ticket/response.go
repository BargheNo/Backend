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
