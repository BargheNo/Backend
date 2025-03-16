package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type CorporationRepository struct{}

func NewCorporationRepository() *CorporationRepository {
	return &CorporationRepository{}
}

func (repo *CorporationRepository) FindCorporationByCIN(db database.Database, cin string) (*entity.Corporation, bool) {
	var corporation entity.Corporation
	result := db.GetDB().Where("cin = ?", cin).First(&corporation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &corporation, true
}

func (repo *CorporationRepository) CreateCorporation(db database.Database, corporation *entity.Corporation) error {
	return db.GetDB().Create(&corporation).Error
}

func (repo *CorporationRepository) DeleteCorporationByCIN(db database.Database, cin string) error {
	return db.GetDB().Where("cin = ?", cin).Delete(&entity.Corporation{}).Error
}

func (repo *CorporationRepository) UpdateCorporation(db database.Database, corporation *entity.Corporation) error {
	return db.GetDB().Save(&corporation).Error
}