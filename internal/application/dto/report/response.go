package reportdto

import (
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	maintenancedto "github.com/BargheNo/Backend/internal/application/dto/maintenance"
)

type MaintenanceReportResponse struct {
	ID                 uint                                           `json:"id"`
	Description        string                                         `json:"description"`
	MaintenanceRequest maintenancedto.AdminMaintenanceRequestResponse `json:"maintenanceRequest"`
	Status             string                                         `json:"status"`
}

type PanelReportResponse struct {
	ID          uint                               `json:"id"`
	Description string                             `json:"description"`
	Panel       installationdto.AdminPanelResponse `json:"panel"`
	Status      string                             `json:"status"`
}
