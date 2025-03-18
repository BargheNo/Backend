package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type BidRepository struct{}

func NewBidRepository() *BidRepository {
	return &BidRepository{}
}

func (repo *BidRepository) GetOpenInstallationRequests(db database.Database, corporationID uint, offset int, pageSize int, sortBy string, dir string) ([]*entity.InstallationRequest, error) {
	var requests []*entity.InstallationRequest
	query := db.GetDB().
		Where("id NOT IN (SELECT request_id FROM bids WHERE corporation_id = ?) AND status = 'open'", corporationID)
	query = query.Order(sortBy + " " + dir)
	result := query.Offset(offset).Limit(pageSize).Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

func (repo *BidRepository) GetRandomOpenInstallationRequests(db database.Database, corporationID uint, offset int, pageSize int) ([]*entity.InstallationRequest, error) {
	var requests []*entity.InstallationRequest
	result := db.GetDB().
		Where("id NOT IN (SELECT request_id FROM bids WHERE corporation_id = ?) AND status = 'open'", corporationID).
		Offset(offset).
		Limit(pageSize).
		Order("RANDOM()").
		Find(&requests)

	if result.Error != nil {
		return nil, result.Error
	}
	return requests, nil
}

func (repo *BidRepository) FindInstallationRequestByID(db database.Database, id uint) (*entity.InstallationRequest, bool) {
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

func (repo *BidRepository) CreateBid(db database.Database, bid *entity.Bid) error {
	return db.GetDB().Create(&bid).Error
}

func (repo *BidRepository) FindBidByID(db database.Database, id uint) (*entity.Bid, bool) {
	var bid entity.Bid
	result := db.GetDB().Where("id = ?", id).First(&bid)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &bid, true
}

func (repo *BidRepository) FindBidByRequestID(db database.Database, requestID uint) (*entity.Bid, bool) {
	var bid entity.Bid
	result := db.GetDB().Where("request_id = ?", requestID).First(&bid)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &bid, true
}

func (repo *BidRepository) DeleteBidByID(db database.Database, id uint) error {
	return db.GetDB().Where("request_id = ?", id).Delete(&entity.Bid{}).Error
}

func (repo *BidRepository) GetBids(db database.Database, corporationID uint, offset int, pageSize int, sortBy string, dir string) ([]*entity.Bid, error) {
	var bids []*entity.Bid
	query := db.GetDB().
		Where("corporation_id = ?", corporationID).
		Order(sortBy + " " + dir)

	result := query.Offset(offset).Limit(pageSize).Find(&bids)
	if result.Error != nil {
		return nil, result.Error
	}
	return bids, nil
}
