package maintenancedto

import (
	"time"

	"github.com/BargheNo/Backend/internal/domain/enum"
)

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

type CorporationMaintenanceListRequest struct {
	CorporationID uint
	OperatorID    uint
	Offset        int
	Limit         int
}

type HandleRequest struct {
	CorporationID uint
	RequestID     uint
	OperatorID    uint
	Accept        bool
}

type AddMaintenanceRecordRequest struct {
	RequestID     uint
	OperatorID    uint
	CorporationID uint
	Date          time.Time
	Title         string
	Details       string
}

type MaintenanceRecordByPanelRequest struct {
	CorporationID uint
	OperatorID    uint
	PanelID       uint
	Offset        int
	Limit         int
}
