package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type GuaranteeRepository interface {
	FindGuaranteeByID(db database.Database, guaranteeID uint) (*entity.Guarantee, bool)
	FindCorporationGuarantee(db database.Database, guaranteeID, corporationID uint) (*entity.Guarantee, bool)
	FindCorporationGuaranteeByName(db database.Database, corporationID uint, name string) (*entity.Guarantee, bool)
	FindCorporationGuarantees(db database.Database, corporationID uint, allowedStatus []enum.GuaranteeStatus) []*entity.Guarantee
	FindGuaranteeTerms(db database.Database, guaranteeID uint) []*entity.GuaranteeTerm
	CreateGuarantee(db database.Database, guarantee *entity.Guarantee) error
	CreateGuaranteeTerms(db database.Database, terms *entity.GuaranteeTerm) error
	UpdateGuarantee(db database.Database, guarantee *entity.Guarantee) error
}
