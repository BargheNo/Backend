package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type BidRepository interface {
	CreateBid(db database.Database, bid *entity.Bid) error
	FindBidByID(db database.Database, id uint) (*entity.Bid, bool)
	FindBidByCorporationAndRequestID(db database.Database, requestID uint, corporationID uint, status []enum.BidStatus) (*entity.Bid, bool)
	DeleteBidByID(db database.Database, id uint) error
	FindCorporationBids(db database.Database, corporationID uint, offset int, pageSize int) []*entity.Bid
	FindRequestBids(db database.Database, requestID uint) []*entity.Bid
}
