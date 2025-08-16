package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
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

func (repo *BidRepository) FindCorporationBids(db database.Database, corporationID uint, allowedStatus []enum.BidStatus, options *postgres.QueryOptions) ([]*entity.Bid, error) {
	var bids []*entity.Bid

	query := db.GetDB().Where("corporation_id = ? AND status IN ?", corporationID, allowedStatus)

	query = applyQueryOptions(query, options)
	result := query.Find(&bids)

	if result.Error != nil {
		return nil, result.Error
	}
	return bids, nil
}

func (repo *BidRepository) CountCorporationBids(db database.Database, corporationID uint, allowedStatus []enum.BidStatus) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Bid{}).
		Where("corporation_id = ? AND status IN ?", corporationID, allowedStatus).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *BidRepository) FindCorporationBidsByQuery(db database.Database, corporationID uint, allowedStatus []enum.BidStatus, query string, options *postgres.QueryOptions) ([]*entity.Bid, error) {
	var bids []*entity.Bid

	result := db.GetDB().
		Joins("LEFT JOIN installation_requests AS requests ON bids.request_id = requests.id").
		Where("corporation_id = ? AND bids.status IN ?", corporationID, allowedStatus).
		Where("name ILIKE ? OR bids.description ILIKE ?", "%"+query+"%", "%"+query+"%")

	result = applyQueryOptions(result, options)
	result = result.Find(&bids)

	if result.Error != nil {
		return nil, result.Error
	}
	return bids, nil
}

func (repo *BidRepository) CountCorporationBidsByQuery(db database.Database, corporationID uint, allowedStatus []enum.BidStatus, query string) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Bid{}).
		Joins("LEFT JOIN installation_requests AS requests ON bids.request_id = requests.id").
		Where("corporation_id = ? AND bids.status IN ?", corporationID, allowedStatus).
		Where("name ILIKE ? OR bids.description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *BidRepository) FindRequestBids(db database.Database, requestID uint, allowedStatus []enum.BidStatus, options *postgres.QueryOptions) ([]*entity.Bid, error) {
	var bids []*entity.Bid

	query := db.GetDB().Where("request_id = ? AND status IN ?", requestID, allowedStatus)
	query = applyQueryOptions(query, options)
	result := query.Find(&bids)

	if result.Error != nil {
		return nil, result.Error
	}
	return bids, nil
}

func (repo *BidRepository) CountRequestBids(
	db database.Database,
	requestID uint,
	allowedStatus []enum.BidStatus,
) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Bid{}).
		Where("request_id = ? AND status IN ?", requestID, allowedStatus).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *BidRepository) FindRequestBidsByQuery(db database.Database, requestID uint, allowedStatus []enum.BidStatus, query string, options *postgres.QueryOptions) ([]*entity.Bid, error) {
	var bids []*entity.Bid

	result := db.GetDB().
		Joins("LEFT JOIN installation_requests AS requests ON bids.request_id = requests.id").
		Where("request_id = ? AND bids.status IN ?", requestID, allowedStatus).
		Where("name ILIKE ? OR bids.description ILIKE ?", "%"+query+"%", "%"+query+"%")

	result = applyQueryOptions(result, options)
	result = result.Find(&bids)

	if result.Error != nil {
		return nil, result.Error
	}
	return bids, nil
}

func (repo *BidRepository) CountRequestBidsByQuery(db database.Database, requestID uint, allowedStatus []enum.BidStatus, query string) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Bid{}).
		Joins("LEFT JOIN installation_requests AS requests ON bids.request_id = requests.id").
		Where("request_id = ? AND bids.status IN ?", requestID, allowedStatus).
		Where("name ILIKE ? OR bids.description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *BidRepository) FindBidsByStatus(db database.Database, allowedStatus []enum.BidStatus, options *postgres.QueryOptions) ([]*entity.Bid, error) {
	var bids []*entity.Bid

	query := db.GetDB().Where("status IN ?", allowedStatus)

	query = applyQueryOptions(query, options)
	result := query.Find(&bids)

	if result.Error != nil {
		return nil, result.Error
	}
	return bids, nil
}

func (repo *BidRepository) CountBidsByStatus(db database.Database, allowedStatus []enum.BidStatus) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Bid{}).
		Where("status IN ?", allowedStatus).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *BidRepository) FindBidsByStatusAndQuery(db database.Database, allowedStatus []enum.BidStatus, query string, options *postgres.QueryOptions) ([]*entity.Bid, error) {
	var bids []*entity.Bid

	result := db.GetDB().
		Joins("LEFT JOIN installation_requests AS requests ON bids.request_id = requests.id").
		Where("bids.status IN ?", allowedStatus).
		Where("name ILIKE ? OR bids.description ILIKE ?", "%"+query+"%", "%"+query+"%")

	result = applyQueryOptions(result, options)
	result = result.Find(&bids)

	if result.Error != nil {
		return nil, result.Error
	}
	return bids, nil
}

func (repo *BidRepository) CountBidsByStatusAndQuery(db database.Database, allowedStatus []enum.BidStatus, query string) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Bid{}).
		Joins("LEFT JOIN installation_requests AS requests ON bids.request_id = requests.id").
		Where("bids.status IN ?", allowedStatus).
		Where("name ILIKE ? OR bids.description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}
