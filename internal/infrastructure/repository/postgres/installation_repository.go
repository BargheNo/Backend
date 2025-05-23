package repositoryimpl

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

func (repo *InstallationRepository) FindRequestByID(db database.Database, requestID uint) (*entity.InstallationRequest, bool) {
	var request entity.InstallationRequest
	result := db.GetDB().First(&request, requestID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &request, true
}

func (repo *InstallationRepository) FindRequestByStatus(db database.Database, status []enum.InstallationRequestStatus, opts ...repository.QueryModifier) []*entity.InstallationRequest {
	var requests []*entity.InstallationRequest
	query := db.GetDB().Where("status IN ?", status)

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}

	result := query.Find(&requests)
	if result.Error != nil {
		panic(result.Error)
	}
	return requests
}

func (repo *InstallationRepository) FindOwnerRequests(db database.Database, ownerID uint, status []enum.InstallationRequestStatus, opts ...repository.QueryModifier) []*entity.InstallationRequest {
	var requests []*entity.InstallationRequest
	query := db.GetDB().Where("owner_id = ? and status IN ?", ownerID, status)

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}

	result := query.Find(&requests)
	if result.Error != nil {
		panic(result.Error)
	}
	return requests
}

func (repo *InstallationRepository) FindOwnerRequestByName(db database.Database, ownerID uint, status []enum.InstallationRequestStatus, name string) (*entity.InstallationRequest, bool) {
	var request *entity.InstallationRequest
	result := db.GetDB().Where("owner_id = ? and name = ? and status IN ?", ownerID, name, status).First(&request)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return request, true
}

func (repo *InstallationRepository) CreateRequest(db database.Database, request *entity.InstallationRequest) error {
	return db.GetDB().Create(&request).Error
}

func (repo *InstallationRepository) CreatePanel(db database.Database, panel *entity.Panel) error {
	return db.GetDB().Create(&panel).Error
}

func (repo *InstallationRepository) FindCorporationPanels(db database.Database, corporationID uint, opts ...repository.QueryModifier) []*entity.Panel {
	var panels []*entity.Panel
	query := db.GetDB().Where("corporation_id = ?", corporationID)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&panels)
	if result.Error != nil {
		panic(result.Error)
	}
	return panels
}

func (repo *InstallationRepository) FindCustomerPanels(db database.Database, customerID uint, opts ...repository.QueryModifier) []*entity.Panel {
	var panels []*entity.Panel
	query := db.GetDB().Where("customer_id = ?", customerID)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&panels)
	if result.Error != nil {
		panic(result.Error)
	}
	return panels
}

func (repo *InstallationRepository) FindPanelByNameAndCustomerID(db database.Database, panelName string, customerID uint) (*entity.Panel, bool) {
	var panel *entity.Panel
	result := db.GetDB().Where("name = ? and customer_id = ?", panelName, customerID).First(&panel)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return panel, true
}

func (repo *InstallationRepository) FindPanelByID(db database.Database, panelID uint) (*entity.Panel, bool) {
	var panel *entity.Panel
	result := db.GetDB().First(&panel, panelID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return panel, true
}

func (repo *InstallationRepository) FindAllPanelsID(db database.Database) []uint {
	var ids []uint
	result := db.GetDB().Model(&entity.Panel{}).Select("id").Find(&ids)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return []uint{}
		}
		panic(result.Error)
	}

	return ids
}
