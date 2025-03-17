package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type CorporationRepository struct{}

func NewCorporationRepository() *CorporationRepository {
	return &CorporationRepository{}
}

func (repo *CorporationRepository) FindCorporationByCIN(db database.Database, cin string) (*entity.Corporation, bool) {
	var corporation entity.Corporation
	result := db.GetDB().Where("cin = ?", cin).First(&corporation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &corporation, true
}

func (repo *CorporationRepository) FindCorporationByID(db database.Database, id uint) (*entity.Corporation, bool) {
	var corporation entity.Corporation
	result := db.GetDB().Where("id = ?", id).First(&corporation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &corporation, true
}

func (repo *CorporationRepository) CreateCorporation(db database.Database, corporation *entity.Corporation) error {
	return db.GetDB().Create(&corporation).Error
}

func (repo *CorporationRepository) DeleteCorporationByCIN(db database.Database, cin string) error {
	return db.GetDB().Where("cin = ?", cin).Delete(&entity.Corporation{}).Error
}

func (repo *CorporationRepository) UpdateCorporation(db database.Database, corporation *entity.Corporation) error {
	return db.GetDB().Save(&corporation).Error
}

func (repo *CorporationRepository) GetOpenInstallationRequests(db database.Database, corporationID uint) ([]*entity.InstallationRequest, error) {
	var requests []*entity.InstallationRequest	
	result := db.GetDB().
		Where("id NOT IN (SELECT request_id FROM bidders WHERE corporation_id = ?) AND status = 'open'", corporationID).
		Find(&requests)
	
	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

func (repo *CorporationRepository) FindInstallationRequestByID(db database.Database, id uint) (*entity.InstallationRequest, bool) {
	var request entity.InstallationRequest
	result := db.GetDB().Where("id = ?", id).First(&request)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &request, true
}

func (repo *CorporationRepository) CreateBidder(db database.Database, bidder *entity.Bidders) error {
	return db.GetDB().Create(&bidder).Error
}

func (repo *CorporationRepository) FindBidderByID(db database.Database, id uint) (*entity.Bidders, bool) {
	var bidder entity.Bidders
	result := db.GetDB().Where("id = ?", id).First(&bidder)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &bidder, true
}

func (repo *CorporationRepository) FindBidderByRequestID(db database.Database, requestID uint) (*entity.Bidders, bool) {
	var bidder entity.Bidders
	result := db.GetDB().Where("request_id = ?", requestID).First(&bidder)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &bidder, true
}

func (repo *CorporationRepository) DeleteBidderByID(db database.Database, id uint) error {
	return db.GetDB().Where("request_id = ?", id).Delete(&entity.Bidders{}).Error
}

func (repo *CorporationRepository) GetBids(db database.Database, corporationID uint) ([]*entity.Bidders, error) {
	var bids []*entity.Bidders
	result := db.GetDB().Where("corporation_id = ?", corporationID).Find(&bids)
	if result.Error != nil {
		return nil, result.Error
	}
	return bids, nil
}