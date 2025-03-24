package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type BidRepository interface {
	CreateBid(db database.Database, bid *entity.Bid) error
	FindBidByID(db database.Database, id uint) (*entity.Bid, bool)
	FindBidByCorporationAndRequestID(db database.Database, requestID uint, corporationID uint) (*entity.Bid, bool)
	DeleteBidByID(db database.Database, id uint) error
	GetBids(db database.Database, corporationID uint, offset int, pageSize int) []*entity.Bid
}
