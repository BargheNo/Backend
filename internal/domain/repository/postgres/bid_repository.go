package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type BidRepository interface {
	CreateBid(db database.Database, bid *entity.Bid) error
	DeleteBidByID(db database.Database, id uint) error
	FindBidByCorporationAndRequestID(db database.Database, requestID uint, corporationID uint, status []enum.BidStatus) (*entity.Bid, bool)
	FindBidByID(db database.Database, id uint) (*entity.Bid, bool)
	FindCorporationBid(db database.Database, bidID uint, corporationID uint) (*entity.Bid, bool)
	FindCorporationBids(db database.Database, corporationID uint, allowedStatus []enum.BidStatus, opts ...QueryModifier) []*entity.Bid
	FindRequestBid(db database.Database, bidID uint, requestID uint) (*entity.Bid, bool)
	FindRequestBids(db database.Database, requestID uint, allowedStatus []enum.BidStatus, opts ...QueryModifier) []*entity.Bid
	UpdateBid(db database.Database, bid *entity.Bid) error
}
