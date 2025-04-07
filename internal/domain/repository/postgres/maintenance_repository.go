package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type MaintenanceRepository interface {
	FindRequestsByPanelID(db database.Database, panelID uint) ([]*entity.MaintenanceRequest, bool)
	CreateMaintenanceRequest(db database.Database, maintenanceRequest *entity.MaintenanceRequest) error
}
