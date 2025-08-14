package usecase

import (
	maintenancedto "github.com/BargheNo/Backend/internal/application/dto/maintenance"
	"github.com/BargheNo/Backend/internal/domain/enum"
)

type MaintenanceService interface {
	GetMaintenanceSortableColumns() []maintenancedto.MaintenanceEnumResponse
	GetMaintenanceUrgencyLevels() []maintenancedto.MaintenanceEnumResponse
	GetMaintenanceRequestStatuses(agent enum.AgentType) []maintenancedto.MaintenanceEnumResponse
	CreateMaintenanceRequest(requestInfo maintenancedto.CreateMaintenanceRequest) error
	GetCustomerMaintenanceRequests(requestInfo maintenancedto.CustomerMaintenanceListRequest) ([]maintenancedto.CustomerMaintenanceRequestResponse, int64, error)
	GetCustomerPanelMaintenanceRequests(listInfo maintenancedto.CustomerPanelMaintenanceListRequest) ([]maintenancedto.CustomerMaintenanceRequestResponse, int64, error)
	GetCustomerMaintenanceRequest(maintenanceInfo maintenancedto.CustomerMaintenanceRequest) (maintenancedto.CustomerMaintenanceRequestResponse, error)
	UpdateMaintenanceRequest(updateRequest maintenancedto.UpdateCustomerRequest) error
	CancelMaintenanceRequest(maintenanceInfo maintenancedto.CustomerMaintenanceRequest) error
	ApproveMaintenanceRecord(maintenanceInfo maintenancedto.CustomerMaintenanceRequest) error
	GetCorporationMaintenanceRequests(listInfo maintenancedto.CorporationMaintenanceListRequest) ([]maintenancedto.CorporationMaintenanceResponse, int64, error)
	GetCorporationMaintenanceRequest(maintenanceInfo maintenancedto.CorporationMaintenanceRequest) (maintenancedto.CorporationMaintenanceResponse, error)
	AcceptMaintenanceRequest(maintenanceInfo maintenancedto.CorporationMaintenanceRequest) error
	RejectMaintenanceRequest(maintenanceInfo maintenancedto.CorporationMaintenanceRequest) error
	CreateMaintenanceRecord(recordInfo maintenancedto.CreateMaintenanceRecordRequest) error
	UpdateMaintenanceRecord(recordInfo maintenancedto.UpdateMaintenanceRecordRequest) error
	ValidateCustomerRecord(recordID, userID uint) error
	GetRequestByAdmin(recordID uint) (maintenancedto.AdminMaintenanceRequestResponse, error)
}
