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
}
