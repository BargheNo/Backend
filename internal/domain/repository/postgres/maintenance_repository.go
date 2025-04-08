package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type MaintenanceRepository interface {
	FindRequestsByPanelID(db database.Database, panelID uint) []*entity.MaintenanceRequest
	CreateMaintenanceRequest(db database.Database, maintenanceRequest *entity.MaintenanceRequest) error
	FindMaintenanceRequestsByOwnerID(db database.Database, ownerID uint, opts ...QueryModifier) []*entity.MaintenanceRequest
	FindMaintenanceRequestsByCorporationID(db database.Database, corporationID uint, opts ...QueryModifier) []*entity.MaintenanceRequest
	FindMaintenanceRequestByID(db database.Database, requestID uint) *entity.MaintenanceRequest
	UpdateMaintenanceRequest(db database.Database, maintenanceRequest *entity.MaintenanceRequest) error
	CreateMaintenanceRecord(db database.Database, maintenanceRecord *entity.MaintenanceRecord) error
	FindMaintenanceRecordsByCorporationID(db database.Database, corporationID uint, opts ...QueryModifier) []*entity.MaintenanceRecord
	FindMaintenanceRecordsByPanelAndCorporationID(db database.Database, panelID uint, corporationID uint, opts ...QueryModifier) []*entity.MaintenanceRecord
}
