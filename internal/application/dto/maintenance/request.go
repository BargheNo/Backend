package maintenancedto

import "github.com/BargheNo/Backend/internal/domain/enum"

type NewMaintenanceRequest struct {
	PanelID       uint              `json:"panelID" validate:"required"`
	OwnerID       uint              `json:"ownerID" validate:"required"`
	CorporationID uint              `json:"corporationID" validate:"required"`
	Subject       string            `json:"subject" validate:"required"`
	Description   string            `json:"description" validate:"required"`
	UrgencyLevel  enum.UrgencyLevel `json:"urgencyLevel" validate:"required"`
}
