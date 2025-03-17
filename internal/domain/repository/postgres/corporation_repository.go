package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type CorporationRepository interface {
	FindCorporationByCIN(db database.Database, cin string) (*entity.Corporation, bool)
	FindCorporationByID(db database.Database, id uint) (*entity.Corporation, bool)
	CreateCorporation(db database.Database, corporation *entity.Corporation) error
	DeleteCorporationByCIN(db database.Database, cin string) error
	UpdateCorporation(db database.Database, corporation *entity.Corporation) error
	GetOpenInstallationRequests(db database.Database, corporationID uint, offset int, pageSize int) ([]*entity.InstallationRequest, error)
	FindInstallationRequestByID(db database.Database, id uint) (*entity.InstallationRequest, bool)
	CreateBid(db database.Database, bid *entity.Bid) error
	FindBidByID(db database.Database, id uint) (*entity.Bid, bool)
	DeleteBidByID(db database.Database, id uint) error
	GetBids(db database.Database, corporationID uint, offset int, pageSize int) ([]*entity.Bid, error)
}
