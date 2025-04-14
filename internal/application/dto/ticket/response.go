package ticketdto

import "time"

type TicketResponse struct {
	ID          uint      `json:"id"`
	OwnerID     uint      `json:"owner_id"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
}
