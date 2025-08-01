package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type InstallationRepository struct{}

func NewInstallationRepository() *InstallationRepository {
	return &InstallationRepository{}
}

func (repo *InstallationRepository) FindRequestByID(db database.Database, requestID uint) (*entity.InstallationRequest, error) {
	var request entity.InstallationRequest
	result := db.GetDB().First(&request, requestID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &request, nil
}

func (repo *InstallationRepository) FindRequestByOwner(db database.Database, requestID, ownerID uint) (*entity.InstallationRequest, error) {
	var request *entity.InstallationRequest
	result := db.GetDB().Where("id = ? AND owner_id = ?", requestID, ownerID).First(&request)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return request, nil
}

func (repo *InstallationRepository) FindRequestsByStatus(db database.Database, status []enum.InstallationRequestStatus, opts ...repository.QueryModifier) ([]*entity.InstallationRequest, error) {
	var requests []*entity.InstallationRequest
	query := db.GetDB().Where("status IN ?", status)

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}

	result := query.Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

func (repo *InstallationRepository) FindOwnerRequests(db database.Database, ownerID uint, status []enum.InstallationRequestStatus, opts ...repository.QueryModifier) ([]*entity.InstallationRequest, error) {
	var requests []*entity.InstallationRequest
	query := db.GetDB().Where("owner_id = ? and status IN ?", ownerID, status)

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}

	result := query.Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

func (repo *InstallationRepository) FindOwnerRequestByName(db database.Database, ownerID uint, status []enum.InstallationRequestStatus, name string) (*entity.InstallationRequest, error) {
	var request *entity.InstallationRequest
	result := db.GetDB().Where("owner_id = ? and name = ? and status IN ?", ownerID, name, status).First(&request)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return request, nil
}

func (repo *InstallationRepository) CreateRequest(db database.Database, request *entity.InstallationRequest) error {
	return db.GetDB().Create(&request).Error
}

func (repo *InstallationRepository) UpdateRequest(db database.Database, request *entity.InstallationRequest) error {
	return db.GetDB().Save(&request).Error
}

func (repo *InstallationRepository) DeleteRequest(db database.Database, request *entity.InstallationRequest) error {
	return db.GetDB().Unscoped().Delete(&request).Error
}

func (repo *InstallationRepository) FindCorporationPanel(db database.Database, panelID, corporationID uint) (*entity.Panel, error) {
	var panel *entity.Panel
	result := db.GetDB().Where("id = ? and corporation_id = ?", panelID, corporationID).First(&panel)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return panel, nil
}

func (repo *InstallationRepository) FindCustomerPanel(db database.Database, panelID, customerID uint) (*entity.Panel, error) {
	var panel *entity.Panel
	result := db.GetDB().Where("id = ? and customer_id = ?", panelID, customerID).First(&panel)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return panel, nil
}

func (repo *InstallationRepository) FindCorporationPanels(db database.Database, corporationID uint, allowedStatus []enum.PanelStatus, opts ...repository.QueryModifier) ([]*entity.Panel, error) {
	var panels []*entity.Panel
	query := db.GetDB().Where("corporation_id = ? AND status IN ?", corporationID, allowedStatus)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&panels)
	if result.Error != nil {
		return nil, result.Error
	}
	return panels, nil
}

func (repo *InstallationRepository) FindPanelsByStatus(db database.Database, allowedStatus []enum.PanelStatus, opts ...repository.QueryModifier) ([]*entity.Panel, error) {
	var panels []*entity.Panel
	query := db.GetDB().Where("status IN ?", allowedStatus)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&panels)
	if result.Error != nil {
		return nil, result.Error
	}
	return panels, nil
}

func (repo *InstallationRepository) FindCustomerPanels(db database.Database, customerID uint, allowedStatus []enum.PanelStatus, opts ...repository.QueryModifier) ([]*entity.Panel, error) {
	var panels []*entity.Panel
	query := db.GetDB().Where("customer_id = ? AND status IN ?", customerID, allowedStatus)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&panels)
	if result.Error != nil {
		return nil, result.Error
	}
	return panels, nil
}

func (repo *InstallationRepository) FindPanelByNameAndCustomerID(db database.Database, panelName string, customerID uint) (*entity.Panel, error) {
	var panel *entity.Panel
	result := db.GetDB().Where("name = ? and customer_id = ?", panelName, customerID).First(&panel)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return panel, nil
}

func (repo *InstallationRepository) FindPanelByID(db database.Database, panelID uint) (*entity.Panel, error) {
	var panel *entity.Panel
	result := db.GetDB().First(&panel, panelID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return panel, nil
}

func (repo *InstallationRepository) FindPanelByOwner(db database.Database, panelID, customerID uint) (*entity.Panel, error) {
	var panel *entity.Panel
	result := db.GetDB().Where("id = ? AND customer_id = ?", panelID, customerID).First(&panel)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return panel, nil
}

func (repo *InstallationRepository) CreatePanel(db database.Database, panel *entity.Panel) error {
	return db.GetDB().Create(&panel).Error
}

func (repo *InstallationRepository) UpdatePanel(db database.Database, panel *entity.Panel) error {
	return db.GetDB().Save(&panel).Error
}

func (repo *InstallationRepository) DeletePanel(db database.Database, panel *entity.Panel) error {
	return db.GetDB().Delete(&panel).Error
}
