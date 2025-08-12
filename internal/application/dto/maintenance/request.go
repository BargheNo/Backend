package maintenancedto

import (
	"time"

	guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"
	"github.com/BargheNo/Backend/internal/domain/enum"
)

type CreateMaintenanceRequest struct {
	PanelID          uint
	OwnerID          uint
	CorporationID    uint
	Subject          string
	Description      string
	UrgencyLevel     enum.UrgencyLevel
	IsUsingGuarantee bool
}

type CustomerMaintenanceListRequest struct {
	Status  uint
	OwnerID uint
	Query   string
	Offset  int
	Limit   int
	SortBy  uint
	Asc     bool
}

type CustomerPanelMaintenanceListRequest struct {
	Status  uint
	PanelID uint
	OwnerID uint
	Offset  int
	Limit   int
	SortBy  uint
	Asc     bool
}

type CustomerMaintenanceRequest struct {
	RequestID uint
	OwnerID   uint
}

type UpdateCustomerRequest struct {
	OwnerID          uint
	RequestID        uint
	Subject          *string
	Description      *string
	UrgencyLevel     *uint
	IsUsingGuarantee *bool
}

type CorporationMaintenanceListRequest struct {
	CorporationID uint
	OperatorID    uint
	Status        uint
	Offset        int
	Limit         int
	SortBy        uint
	Asc           bool
}

type CorporationMaintenanceRequest struct {
	CorporationID uint
	OperatorID    uint
	RequestID     uint
}

type CreateMaintenanceRecordRequest struct {
	CorporationID      uint
	OperatorID         uint
	RequestID          uint
	Title              string
	Details            string
	GuaranteeViolation *guaranteedto.CreateGuaranteeViolationRequest
}

type UpdateMaintenanceRecordRequest struct {
	CorporationID      uint
	OperatorID         uint
	RequestID          uint
	Title              *string
	Details            *string
	GuaranteeViolation *guaranteedto.UpdateGuaranteeViolationRequest
}

type HandleRequest struct {
	CorporationID uint
	RequestID     uint
	OperatorID    uint
	Accept        bool
}

type AddMaintenanceRecordRequest struct {
	OperatorID    uint
	CorporationID uint
	PanelID       uint
	Date          time.Time
	Title         string
	Details       string
}

type CorporationMaintenanceRecordByPanelRequest struct {
	CorporationID uint
	OperatorID    uint
	PanelID       uint
	Offset        int
	Limit         int
}

type CustomerMaintenanceRecordByPanelRequest struct {
	OwnerID uint
	PanelID uint
	Offset  int
	Limit   int
}
