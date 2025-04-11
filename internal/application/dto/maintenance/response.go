package maintenancedto

import (
	"time"

	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
)

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
	Panel         installationdto.CustomerPanelResponse
}

type CorporationMaintenanceResponse struct {
	ID           uint
	PanelID      uint
	Subject      string
	Description  string
	UrgencyLevel string
	Status       string
	CreatedAt    time.Time
	OwnerPhone   string
	Panel        installationdto.CorporationPanelResponse
}

type MaintenanceRecordResponse struct {
	ID            uint
	RequestID     uint
	Panel         installationdto.CorporationPanelResponse
	OperatorID    uint
	CorporationID uint
	Title         string
	Details       string
	Date          time.Time
}

type CustomerMaintenanceRecordResponse struct {
	ID            uint
	Panel         installationdto.CustomerPanelResponse
	OperatorID    uint
	OperatorPhone string
	Title         string
	Details       string
	Date          time.Time
}
