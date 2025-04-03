package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type BidRepository struct{}

func NewBidRepository() *BidRepository {
	return &BidRepository{}
}

func (repo *BidRepository) CreateBid(db database.Database, bid *entity.Bid) error {
	return db.GetDB().Create(&bid).Error
}

func (repo *BidRepository) FindBidByID(db database.Database, id uint) (*entity.Bid, bool) {
	var bid entity.Bid
	result := db.GetDB().First(&bid, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &bid, true
}

func (repo *BidRepository) FindBidByCorporationAndRequestID(db database.Database, requestID uint, corporationID uint, status []enum.BidStatus) (*entity.Bid, bool) {
	var bid entity.Bid
	result := db.GetDB().Where("request_id = ? AND corporation_id = ? AND status IN ?", requestID, corporationID, status).First(&bid)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &bid, true
}

func (repo *BidRepository) DeleteBidByID(db database.Database, id uint) error {
	return db.GetDB().Where("id = ?", id).Delete(&entity.Bid{}).Error
}

func (repo *BidRepository) FindCorporationBids(db database.Database, corporationID uint, offset int, pageSize int) []*entity.Bid {
	var bids []*entity.Bid

	result := db.GetDB().Model(&entity.Bid{}).
		Where("corporation_id = ?", corporationID).
		Offset(offset).
		Limit(pageSize).
		Find(&bids)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return bids
}

func (repo *BidRepository) FindRequestBids(db database.Database, requestID uint) []*entity.Bid {
	var bids []*entity.Bid
	result := db.GetDB().Model(&entity.Bid{}).Where("request_id = ?", requestID).Find(&bids)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return bids
}
