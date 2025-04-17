package reportdto

import maintenancedto "github.com/BargheNo/Backend/internal/application/dto/maintenance"

type MaintenanceReportResponse struct {
	ID                uint
	Description       string
	MaintenanceRecord maintenancedto.MaintenanceRecordResponse
	Status            string
}
