package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type BidRepository struct{}

func NewBidRepository() *BidRepository {
	return &BidRepository{}
}

func (repo *BidRepository) FindBidByID(db database.Database, id uint) (*entity.Bid, error) {
	var bid entity.Bid
	result := db.GetDB().First(&bid, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &bid, nil
}

func (repo *BidRepository) FindRequestBid(db database.Database, bidID, requestID uint) (*entity.Bid, error) {
	var bid entity.Bid
	result := db.GetDB().Where("id = ? AND request_id = ?", bidID, requestID).First(&bid)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &bid, nil
}

func (repo *BidRepository) FindCorporationBid(db database.Database, bidID, corporationID uint) (*entity.Bid, error) {
	var bid entity.Bid
	result := db.GetDB().Where("id = ? AND corporation_id = ?", bidID, corporationID).First(&bid)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &bid, nil
}

func (repo *BidRepository) FindBidByCorporationAndRequestID(db database.Database, requestID uint, corporationID uint, status []enum.BidStatus) (*entity.Bid, error) {
	var bid entity.Bid
	result := db.GetDB().Where("request_id = ? AND corporation_id = ? AND status IN ?", requestID, corporationID, status).First(&bid)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &bid, nil
}

func (repo *BidRepository) DeleteBidByID(db database.Database, id uint) error {
	return db.GetDB().Where("id = ?", id).Delete(&entity.Bid{}).Error
}

func (repo *BidRepository) CreateBid(db database.Database, bid *entity.Bid) error {
	return db.GetDB().Create(&bid).Error
}

func (repo *BidRepository) UpdateBid(db database.Database, bid *entity.Bid) error {
	return db.GetDB().Save(&bid).Error
}

func (repo *BidRepository) FindCorporationBids(db database.Database, corporationID uint, allowedStatus []enum.BidStatus, opts ...repository.QueryModifier) ([]*entity.Bid, error) {
	var bids []*entity.Bid

	query := db.GetDB().Where("corporation_id = ? AND status IN ?", corporationID, allowedStatus)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&bids)

	if result.Error != nil {
		return nil, result.Error
	}
	return bids, nil
}

func (repo *BidRepository) FindRequestBids(db database.Database, requestID uint, allowedStatus []enum.BidStatus, opts ...repository.QueryModifier) ([]*entity.Bid, error) {
	var bids []*entity.Bid

	query := db.GetDB().Where("request_id = ? AND status IN ?", requestID, allowedStatus)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&bids)

	if result.Error != nil {
		return nil, result.Error
	}
	return bids, nil
}
