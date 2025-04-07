package service

import maintenancedto "github.com/BargheNo/Backend/internal/application/dto/maintenance"

type MaintenanceService interface {
	CreateMaintenanceRequest(requestInfo maintenancedto.NewMaintenanceRequest)
}
