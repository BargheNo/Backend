package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type MaintenanceRepository struct{}

func NewMaintenanceRepository() *MaintenanceRepository {
	return &MaintenanceRepository{}
}

func (repo *MaintenanceRepository) FindRequestByID(db database.Database, requestID uint) (*entity.MaintenanceRequest, error) {
	var request *entity.MaintenanceRequest
	result := db.GetDB().First(&request, requestID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return request, nil
}

func (repo *MaintenanceRepository) FindRequestsByPanelIDAndStatus(db database.Database, panelID uint, allowedStatus []enum.MaintenanceRequestStatus, opts ...repository.QueryModifier) ([]*entity.MaintenanceRequest, error) {
	var requests []*entity.MaintenanceRequest
	query := db.GetDB().Where("panel_id = ? AND status IN ?", panelID, allowedStatus)

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}

	result := query.Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

func (repo *MaintenanceRepository) FindRequestsByCustomerID(db database.Database, customerID uint, allowedStatus []enum.MaintenanceRequestStatus, opts ...repository.QueryModifier) ([]*entity.MaintenanceRequest, error) {
	var requests []*entity.MaintenanceRequest
	query := db.GetDB().
		Joins("LEFT JOIN panels AS Panel ON maintenance_requests.panel_id = Panel.id").
		Where("Panel.customer_id = ?", customerID).
		Where("maintenance_requests.status IN ?", allowedStatus)

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}

	result := query.Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

func (repo *MaintenanceRepository) CreateMaintenanceRequest(db database.Database, maintenanceRequest *entity.MaintenanceRequest) error {
	return db.GetDB().Create(maintenanceRequest).Error
}

func (repo *MaintenanceRepository) FindMaintenanceRequestsByOwnerID(db database.Database, ownerID uint, opts ...repository.QueryModifier) ([]*entity.MaintenanceRequest, error) {
	var requests []*entity.MaintenanceRequest
	query := db.GetDB().Where("owner_id = ?", ownerID)

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}

	result := query.Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

func (repo *MaintenanceRepository) FindCorporationRequestsByStatus(db database.Database, corporationID uint, allowedStatuses []enum.MaintenanceRequestStatus, opts ...repository.QueryModifier) ([]*entity.MaintenanceRequest, error) {
	var requests []*entity.MaintenanceRequest
	query := db.GetDB().Where("corporation_id = ? AND  status IN ?", corporationID, allowedStatuses)

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}

	result := query.Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

func (repo *MaintenanceRepository) FindCorporationRequestByStatus(db database.Database, requestID, corporationID uint, allowedStatus []enum.MaintenanceRequestStatus) (*entity.MaintenanceRequest, error) {
	var request *entity.MaintenanceRequest
	result := db.GetDB().Where("id = ? AND corporation_id = ? AND status IN ?", requestID, corporationID, allowedStatus).First(&request)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return request, nil
}

func (repo *MaintenanceRepository) FindRecordByRequestID(db database.Database, requestID uint) (*entity.MaintenanceRecord, error) {
	var record *entity.MaintenanceRecord
	result := db.GetDB().Where("request_id = ?", requestID).First(&record)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return record, nil
}

func (repo *MaintenanceRepository) FindRecordByID(db database.Database, recordID uint) (*entity.MaintenanceRecord, error) {
	var record *entity.MaintenanceRecord
	result := db.GetDB().First(&record, recordID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return record, nil
}

func (repo *MaintenanceRepository) UpdateMaintenanceRequest(db database.Database, maintenanceRequest *entity.MaintenanceRequest) error {
	return db.GetDB().Save(maintenanceRequest).Error
}

func (repo *MaintenanceRepository) CreateMaintenanceRecord(db database.Database, maintenanceRecord *entity.MaintenanceRecord) error {
	return db.GetDB().Create(maintenanceRecord).Error
}

func (repo *MaintenanceRepository) UpdateMaintenanceRecord(db database.Database, maintenanceRecord *entity.MaintenanceRecord) error {
	return db.GetDB().Save(maintenanceRecord).Error
}

func (repo *MaintenanceRepository) FindMaintenanceRecordsByPanelAndCorporationID(db database.Database, panelID uint, corporationID uint, opts ...repository.QueryModifier) ([]*entity.MaintenanceRecord, error) {
	var records []*entity.MaintenanceRecord
	query := db.GetDB().Where("panel_id = ? AND corporation_id = ?", panelID, corporationID)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}
	return records, nil
}
