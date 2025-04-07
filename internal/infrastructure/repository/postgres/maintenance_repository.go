package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type MaintenanceRepository struct{}

func NewMaintenanceRepository() *MaintenanceRepository {
	return &MaintenanceRepository{}
}

func (repo *MaintenanceRepository) FindRequestsByPanelID(db database.Database, panelID uint) []*entity.MaintenanceRequest {
	var requests []*entity.MaintenanceRequest
	if err := db.GetDB().Where("panel_id = ?", panelID).Find(&requests).Error; err != nil {
		return nil
	}
	if len(requests) == 0 {
		return nil
	}
	return requests
}

func (repo *MaintenanceRepository) CreateMaintenanceRequest(db database.Database, maintenanceRequest *entity.MaintenanceRequest) error {
	return db.GetDB().Create(maintenanceRequest).Error
}

func (repo *MaintenanceRepository) FindMaintenanceRequestsByOwnerID(db database.Database, ownerID uint, opts ...repository.QueryModifier) []*entity.MaintenanceRequest {
	var requests []*entity.MaintenanceRequest
	query := db.GetDB().Where("owner_id = ?", ownerID)

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}

	result := query.Find(&requests)
	if result.Error != nil {
		panic(result.Error)
	}
	return requests
}
