package reportdto

import (
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	maintenancedto "github.com/BargheNo/Backend/internal/application/dto/maintenance"
)

type MaintenanceReportResponse struct {
	ID                uint                                     `json:"id"`
	Description       string                                   `json:"description"`
	MaintenanceRecord maintenancedto.MaintenanceRecordResponse `json:"maintenanceRecord"`
	Status            string                                   `json:"status"`
}

type PanelReportResponse struct {
	ID          uint                          `json:"id"`
	Description string                        `json:"description"`
	Panel       installationdto.PanelResponse `json:"panel"`
	Status      string                        `json:"status"`
}
