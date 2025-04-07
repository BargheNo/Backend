package service

import maintenancedto "github.com/BargheNo/Backend/internal/application/dto/maintenance"

type MaintenanceService interface {
	CreateMaintenanceRequest(requestInfo maintenancedto.NewMaintenanceRequest)
	GetCustomerMaintenanceRequests(requestInfo maintenancedto.MaintenanceListRequest) []maintenancedto.MaintenanceResponse
	GetCorporationMaintenanceRequests(requestInfo maintenancedto.CorporationMaintenanceListRequest) []maintenancedto.CorporationMaintenanceResponse
}
