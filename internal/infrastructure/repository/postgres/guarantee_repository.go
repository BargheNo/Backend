package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type GuaranteeRepository struct{}

func NewGuaranteeRepository() *GuaranteeRepository {
	return &GuaranteeRepository{}
}

func (repo *GuaranteeRepository) FindGuaranteeByID(db database.Database, guaranteeID uint) (*entity.Guarantee, bool) {
	var guarantee *entity.Guarantee
	result := db.GetDB().First(&guarantee, guaranteeID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return guarantee, true
}

func (repo *GuaranteeRepository) FindCorporationGuaranteeByName(db database.Database, corporationID uint, name string) (*entity.Guarantee, bool) {
	var guarantee *entity.Guarantee
	result := db.GetDB().Where("corporation_id = ? AND name = ?", corporationID, name).First(&guarantee)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return guarantee, true
}

func (repo *GuaranteeRepository) FindCorporationGuarantee(db database.Database, guaranteeID, corporationID uint) (*entity.Guarantee, bool) {
	var guarantee *entity.Guarantee
	result := db.GetDB().Where("id = ? AND corporation_id = ?", guaranteeID, corporationID).First(&guarantee)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return guarantee, true
}

func (repo *GuaranteeRepository) FindCorporationGuarantees(db database.Database, corporationID uint, allowedStatus []enum.GuaranteeStatus) []*entity.Guarantee {
	var guarantees []*entity.Guarantee
	result := db.GetDB().Where("corporation_id = ? AND status IN ?", corporationID, allowedStatus).Find(&guarantees)
	if result.Error != nil {
		panic(result.Error)
	}
	return guarantees
}

func (repo *GuaranteeRepository) FindGuaranteeTerms(db database.Database, guaranteeID uint) []*entity.GuaranteeTerm {
	var terms []*entity.GuaranteeTerm
	result := db.GetDB().Where("guarantee_id = ?", guaranteeID).Find(&terms)
	if result.Error != nil {
		panic(result.Error)
	}
	return terms
}

func (repo *GuaranteeRepository) CreateGuarantee(db database.Database, guarantee *entity.Guarantee) error {
	return db.GetDB().Create(&guarantee).Error
}

func (repo *GuaranteeRepository) CreateGuaranteeTerms(db database.Database, terms *entity.GuaranteeTerm) error {
	return db.GetDB().Create(&terms).Error
}

func (repo *GuaranteeRepository) UpdateGuarantee(db database.Database, guarantee *entity.Guarantee) error {
	return db.GetDB().Save(&guarantee).Error
}

func (repo *GuaranteeRepository) FindPanelGuaranteeViolation(db database.Database, panelID uint) (*entity.GuaranteeViolation, bool) {
	var violation *entity.GuaranteeViolation
	result := db.GetDB().Where("panel_id = ?", panelID).First(&violation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return violation, true
}

func (repo *GuaranteeRepository) CreateGuaranteeViolation(db database.Database, violation *entity.GuaranteeViolation) error {
	return db.GetDB().Create(violation).Error
}

func (repo *GuaranteeRepository) UpdateGuaranteeViolation(db database.Database, violation *entity.GuaranteeViolation) error {
	return db.GetDB().Save(violation).Error
}

func (repo *GuaranteeRepository) DeleteGuaranteeViolation(db database.Database, violation *entity.GuaranteeViolation) error {
	return db.GetDB().Unscoped().Delete(violation).Error
}
