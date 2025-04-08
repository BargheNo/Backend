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

func (repo *MaintenanceRepository) FindMaintenanceRequestsByCorporationID(db database.Database, corporationID uint, opts ...repository.QueryModifier) []*entity.MaintenanceRequest {
	var requests []*entity.MaintenanceRequest
	query := db.GetDB().Where("corporation_id = ?", corporationID)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&requests)
	if result.Error != nil {
		panic(result.Error)
	}

	return requests
}

func (repo *MaintenanceRepository) FindMaintenanceRequestByID(db database.Database, requestID uint) *entity.MaintenanceRequest {
	var request entity.MaintenanceRequest
	if err := db.GetDB().First(&request, requestID).Error; err != nil {
		return nil
	}
	return &request
}

func (repo *MaintenanceRepository) UpdateMaintenanceRequest(db database.Database, maintenanceRequest *entity.MaintenanceRequest) error {
	return db.GetDB().Save(maintenanceRequest).Error
}

func (repo *MaintenanceRepository) CreateMaintenanceRecord(db database.Database, maintenanceRecord *entity.MaintenanceRecord) error {
	return db.GetDB().Create(maintenanceRecord).Error
}

func (repo *MaintenanceRepository) FindMaintenanceRecordsByCorporationID(db database.Database, corporationID uint, opts ...repository.QueryModifier) []*entity.MaintenanceRecord {
	var records []*entity.MaintenanceRecord
	query := db.GetDB().Where("corporation_id = ?", corporationID)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&records)
	if result.Error != nil {
		panic(result.Error)
	}

	return records
}

func (repo *MaintenanceRepository) FindMaintenanceRecordsByPanelAndCorporationID(db database.Database, panelID uint, corporationID uint, opts ...repository.QueryModifier) []*entity.MaintenanceRecord {
	var records []*entity.MaintenanceRecord
	query := db.GetDB().Where("panel_id = ? AND corporation_id = ?", panelID, corporationID)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&records)
	if result.Error != nil {
		panic(result.Error)
	}
	return records
}

func (repo *MaintenanceRepository) FindMaintenanceRecordsByCustomerID(db database.Database, customerID uint, opts ...repository.QueryModifier) []*entity.MaintenanceRecord {
	var records []*entity.MaintenanceRecord
	query := db.GetDB().Where("customer_id = ?", customerID)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&records)
	if result.Error != nil {
		panic(result.Error)
	}
	return records
}

func (repo *MaintenanceRepository) FindCustomerMaintenanceRecordsByPanelID(db database.Database, customerID uint, panelID uint, opts ...repository.QueryModifier) []*entity.MaintenanceRecord {
	var records []*entity.MaintenanceRecord
	print("customerID", customerID)
	print("panelID", panelID)
	query := db.GetDB().Where("customer_id = ? AND panel_id = ?", customerID, panelID)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&records)
	if result.Error != nil {
		panic(result.Error)
	}
	return records
}
