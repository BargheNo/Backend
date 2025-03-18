package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type BidRepository interface {
	GetOpenInstallationRequests(db database.Database, corporationID uint, offset int, pageSize int, sortBy string, dir string) ([]*entity.InstallationRequest, error)
	GetRandomOpenInstallationRequests(db database.Database, corporationID uint, pageSize int, offset int) ([]*entity.InstallationRequest, error)
	FindInstallationRequestByID(db database.Database, id uint) (*entity.InstallationRequest, bool)
	CreateBid(db database.Database, bid *entity.Bid) error
	FindBidByID(db database.Database, id uint) (*entity.Bid, bool)
	DeleteBidByID(db database.Database, id uint) error
	GetBids(db database.Database, corporationID uint, offset int, pageSize int, sortBy string, dir string) ([]*entity.Bid, error)
}
