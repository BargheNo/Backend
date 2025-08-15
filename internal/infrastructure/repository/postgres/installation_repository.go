package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type InstallationRepository struct{}

func NewInstallationRepository() *InstallationRepository {
	return &InstallationRepository{}
}

const (
	queryByStatus = "status IN ?"
)

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

func (repo *InstallationRepository) FindRequestsByStatus(db database.Database, status []enum.InstallationRequestStatus, options *postgres.QueryOptions) ([]*entity.InstallationRequest, error) {
	var requests []*entity.InstallationRequest
	query := db.GetDB().Where(queryByStatus, status)

	query = applyQueryOptions(query, options)

	result := query.Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

func (repo *InstallationRepository) CountRequestsByStatus(db database.Database, status []enum.InstallationRequestStatus) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.InstallationRequest{}).
		Where(queryByStatus, status).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *InstallationRepository) FindCorporationRequestsByQuery(db database.Database, allowedStatus []enum.InstallationRequestStatus, query string, options *postgres.QueryOptions) ([]*entity.InstallationRequest, error) {
	var requests []*entity.InstallationRequest
	result := db.GetDB().
		Where(queryByStatus, allowedStatus).
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")

	result = applyQueryOptions(result, options)
	result = result.Find(&requests)

	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

func (repo *InstallationRepository) CountCorporationRequestsByQuery(db database.Database, allowedStatus []enum.InstallationRequestStatus, query string) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Panel{}).
		Where(queryByStatus, allowedStatus).
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *InstallationRepository) FindRequestsByQuery(db database.Database, query string, options *postgres.QueryOptions) ([]*entity.InstallationRequest, error) {
	var requests []*entity.InstallationRequest
	result := db.GetDB().
		Joins("LEFT JOIN users AS customers ON installation_requests.owner_id = customers.id").
		Where(`
			installation_requests.name ILIKE ? OR 
			installation_requests.description ILIKE ? OR
			customers.first_name ILIKE ? OR 
			customers.last_name ILIKE ? OR 
			customers.email ILIKE ? OR
			customers.phone ILIKE ?
		`,
			"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%",
			"%"+query+"%")

	result = applyQueryOptions(result, options)

	result = result.Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

func (repo *InstallationRepository) CountRequestsByQuery(db database.Database, query string) (int64, error) {
	var count int64
	err := db.GetDB().
		Model(&entity.InstallationRequest{}).
		Joins("LEFT JOIN users AS customers ON installation_requests.owner_id = customers.id").
		Where(`
			installation_requests.name ILIKE ? OR 
			installation_requests.description ILIKE ? OR
			customers.first_name ILIKE ? OR 
			customers.last_name ILIKE ? OR 
			customers.email ILIKE ? OR
			customers.phone ILIKE ?
		`,
			"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%",
			"%"+query+"%").
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *InstallationRepository) FindOwnerRequests(db database.Database, ownerID uint, status []enum.InstallationRequestStatus, options *postgres.QueryOptions) ([]*entity.InstallationRequest, error) {
	var requests []*entity.InstallationRequest
	query := db.GetDB().Where("owner_id = ? and status IN ?", ownerID, status)

	query = applyQueryOptions(query, options)

	result := query.Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

func (repo *InstallationRepository) CountOwnerRequests(db database.Database, ownerID uint, status []enum.InstallationRequestStatus) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.InstallationRequest{}).
		Where("owner_id = ? and status IN ?", ownerID, status).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
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

func (repo *InstallationRepository) FindCorporationPanels(db database.Database, corporationID uint, allowedStatus []enum.PanelStatus, options *postgres.QueryOptions) ([]*entity.Panel, error) {
	var panels []*entity.Panel
	query := db.GetDB().Where("corporation_id = ? AND status IN ?", corporationID, allowedStatus)
	query = applyQueryOptions(query, options)
	result := query.Find(&panels)
	if result.Error != nil {
		return nil, result.Error
	}
	return panels, nil
}

func (repo *InstallationRepository) CountCorporationPanels(db database.Database, corporationID uint, allowedStatus []enum.PanelStatus) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Panel{}).
		Where("corporation_id = ? AND status IN ?", corporationID, allowedStatus).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *InstallationRepository) FindPanelsByStatus(db database.Database, allowedStatus []enum.PanelStatus, options *postgres.QueryOptions) ([]*entity.Panel, error) {
	var panels []*entity.Panel
	query := db.GetDB().Where(queryByStatus, allowedStatus)
	query = applyQueryOptions(query, options)
	result := query.Find(&panels)
	if result.Error != nil {
		return nil, result.Error
	}
	return panels, nil
}

func (repo *InstallationRepository) CountPanelsByStatus(db database.Database, allowedStatus []enum.PanelStatus) (int64, error) {
	var count int64

	err := db.GetDB().Model(&entity.Panel{}).Where(queryByStatus, allowedStatus).Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *InstallationRepository) FindPanelsByQuery(db database.Database, allowedStatus []enum.PanelStatus, query string, options *postgres.QueryOptions) ([]*entity.Panel, error) {
	var panels []*entity.Panel

	result := db.GetDB().
		Joins("LEFT JOIN users AS customers ON panels.customer_id = customers.id").
		Joins("LEFT JOIN corporations ON panels.corporation_id = corporations.id").
		Where(`
			panels.status IN ? AND (
			panels.name ILIKE ? OR 
			customers.first_name ILIKE ? OR 
			customers.last_name ILIKE ? OR 
			customers.email ILIKE ? OR
			customers.phone ILIKE ? OR
			corporations.name ILIKE ? OR
			corporations.registration_number ILIKE ?
		)
		`,
			allowedStatus, "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%",
			"%"+query+"%", "%"+query+"%", "%"+query+"%")

	result = applyQueryOptions(result, options)
	result = result.Find(&panels)

	if result.Error != nil {
		return nil, result.Error
	}
	return panels, nil
}

func (repo *InstallationRepository) CountPanelsByQuery(db database.Database, allowedStatus []enum.PanelStatus, query string) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Panel{}).
		Joins("LEFT JOIN users AS customers ON panels.customer_id = customers.id").
		Joins("LEFT JOIN corporations ON panels.corporation_id = corporations.id").
		Where(`
			panels.status IN ? AND (
			panels.name ILIKE ? OR 
			customers.first_name ILIKE ? OR 
			customers.last_name ILIKE ? OR 
			customers.email ILIKE ? OR
			customers.phone ILIKE ? OR
			corporations.name ILIKE ? OR
			corporations.registration_number ILIKE ?
		)
		`,
			allowedStatus, "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%",
			"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *InstallationRepository) FindCustomerPanels(db database.Database, customerID uint, allowedStatus []enum.PanelStatus, options *postgres.QueryOptions) ([]*entity.Panel, error) {
	var panels []*entity.Panel
	query := db.GetDB().Where("customer_id = ? AND status IN ?", customerID, allowedStatus)
	query = applyQueryOptions(query, options)
	result := query.Find(&panels)
	if result.Error != nil {
		return nil, result.Error
	}
	return panels, nil
}

func (repo *InstallationRepository) CountCustomerPanels(db database.Database, customerID uint, allowedStatus []enum.PanelStatus) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Panel{}).
		Where("customer_id = ? AND status IN ?", customerID, allowedStatus).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *InstallationRepository) FindCustomerPanelsByQuery(db database.Database, customerID uint, allowedStatus []enum.PanelStatus, query string, options *postgres.QueryOptions) ([]*entity.Panel, error) {
	var panels []*entity.Panel
	result := db.GetDB().
		Where("customer_id = ? AND status IN ?", customerID, allowedStatus).
		Where("name ILIKE ?", "%"+query+"%")

	result = applyQueryOptions(result, options)
	result = result.Find(&panels)

	if result.Error != nil {
		return nil, result.Error
	}
	return panels, nil
}

func (repo *InstallationRepository) CountCustomerPanelsByQuery(db database.Database, customerID uint, allowedStatus []enum.PanelStatus, query string) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Panel{}).
		Where("customer_id = ? AND status IN ?", customerID, allowedStatus).
		Where("name ILIKE ?", "%"+query+"%").
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *InstallationRepository) FindCorporationPanelsByQuery(db database.Database, corporationID uint, allowedStatus []enum.PanelStatus, query string, options *postgres.QueryOptions) ([]*entity.Panel, error) {
	var panels []*entity.Panel
	result := db.GetDB().
		Where("corporation_id = ? AND status IN ?", corporationID, allowedStatus).
		Where("name ILIKE ?", "%"+query+"%")

	result = applyQueryOptions(result, options)
	result = result.Find(&panels)

	if result.Error != nil {
		return nil, result.Error
	}
	return panels, nil
}

func (repo *InstallationRepository) CountCorporationPanelsByQuery(db database.Database, corporationID uint, allowedStatus []enum.PanelStatus, query string) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Panel{}).
		Where("corporation_id = ? AND status IN ?", corporationID, allowedStatus).
		Where("name ILIKE ?", "%"+query+"%").
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
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
	return db.GetDB().Unscoped().Delete(&panel).Error
}
