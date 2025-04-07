package maintenancedto

import "time"

type MaintenanceResponse struct {
	ID            uint
	PanelID       uint
	CorporationID uint
	OwnerID       uint
	Subject       string
	Description   string
	UrgencyLevel  string
	Status        string
	CreatedAt     time.Time
}
