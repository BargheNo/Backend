package reportdto

import (
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	maintenancedto "github.com/BargheNo/Backend/internal/application/dto/maintenance"
)

type MaintenanceReportResponse struct {
	ID                uint
	Description       string
	MaintenanceRecord maintenancedto.MaintenanceRecordResponse
	Status            string
}

type PanelReportResponse struct {
	ID          uint
	Description string
	Panel       installationdto.PanleResponse
	Status      string
}
