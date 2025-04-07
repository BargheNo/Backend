package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type MaintenanceRepository struct{}

func NewMaintenanceRepository() *MaintenanceRepository {
	return &MaintenanceRepository{}
}

func (repo *MaintenanceRepository) FindRequestsByPanelID(db database.Database, panelID uint) ([]*entity.MaintenanceRequest, bool) {
	var requests []*entity.MaintenanceRequest
	if err := db.GetDB().Where("panel_id = ?", panelID).Find(&requests).Error; err != nil {
		return nil, false
	}
	if len(requests) == 0 {
		return nil, false
	}
	return requests, true
}

func (repo *MaintenanceRepository) CreateMaintenanceRequest(db database.Database, maintenanceRequest *entity.MaintenanceRequest) error {
	return db.GetDB().Create(maintenanceRequest).Error
}
