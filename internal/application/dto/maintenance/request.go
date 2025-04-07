package maintenancedto

import "github.com/BargheNo/Backend/internal/domain/enum"

type NewMaintenanceRequest struct {
	PanelID       uint
	OwnerID       uint
	CorporationID uint
	Subject       string
	Description   string
	UrgencyLevel  enum.UrgencyLevel
}

type MaintenanceListRequest struct {
	OwnerID uint
	Offset  int
	Limit   int
}
