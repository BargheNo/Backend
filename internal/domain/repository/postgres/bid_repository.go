package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type BidRepository interface {
	CreateBid(db database.Database, bid *entity.Bid) error
	DeleteBidByID(db database.Database, id uint) error
	FindBidByCorporationAndRequestID(db database.Database, requestID uint, corporationID uint, status []enum.BidStatus) (*entity.Bid, error)
	FindBidByID(db database.Database, id uint) (*entity.Bid, error)
	FindCorporationBid(db database.Database, bidID uint, corporationID uint) (*entity.Bid, error)
	FindCorporationBids(db database.Database, corporationID uint, allowedStatus []enum.BidStatus, options *QueryOptions) ([]*entity.Bid, error)
	CountCorporationBids(db database.Database, corporationID uint, allowedStatus []enum.BidStatus) (int64, error)
	FindCorporationBidsByQuery(db database.Database, corporationID uint, allowedStatus []enum.BidStatus, query string, options *QueryOptions) ([]*entity.Bid, error)
	CountCorporationBidsByQuery(db database.Database, corporationID uint, allowedStatus []enum.BidStatus, query string) (int64, error)
	FindRequestBid(db database.Database, bidID uint, requestID uint) (*entity.Bid, error)
	FindRequestBids(db database.Database, requestID uint, allowedStatus []enum.BidStatus, options *QueryOptions) ([]*entity.Bid, error)
	CountRequestBids(db database.Database, requestID uint, allowedStatus []enum.BidStatus) (int64, error)
	FindRequestBidsByQuery(db database.Database, requestID uint, allowedStatus []enum.BidStatus, query string, options *QueryOptions) ([]*entity.Bid, error)
	CountRequestBidsByQuery(db database.Database, requestID uint, allowedStatus []enum.BidStatus, query string) (int64, error)
	FindBidsByStatus(db database.Database, allowedStatus []enum.BidStatus, options *QueryOptions) ([]*entity.Bid, error)
	CountBidsByStatus(db database.Database, allowedStatus []enum.BidStatus) (int64, error)
	FindBidsByStatusAndQuery(db database.Database, allowedStatus []enum.BidStatus, query string, options *QueryOptions) ([]*entity.Bid, error)
	CountBidsByStatusAndQuery(db database.Database, allowedStatus []enum.BidStatus, query string) (int64, error)
	UpdateBid(db database.Database, bid *entity.Bid) error
}
