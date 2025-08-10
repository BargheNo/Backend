package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type GuaranteeRepository interface {
	FindGuaranteeByID(db database.Database, guaranteeID uint) (*entity.Guarantee, error)
	FindCorporationGuarantee(db database.Database, guaranteeID, corporationID uint) (*entity.Guarantee, error)
	FindCorporationGuaranteeByName(db database.Database, corporationID uint, name string) (*entity.Guarantee, error)
	FindCorporationGuarantees(db database.Database, corporationID uint, allowedStatus []enum.GuaranteeStatus) ([]*entity.Guarantee, error)
	FindGuaranteeTerms(db database.Database, guaranteeID uint) ([]*entity.GuaranteeTerm, error)
	CreateGuarantee(db database.Database, guarantee *entity.Guarantee) error
	CreateGuaranteeTerms(db database.Database, terms *entity.GuaranteeTerm) error
	UpdateGuarantee(db database.Database, guarantee *entity.Guarantee) error
	FindPanelGuaranteeViolation(db database.Database, panelID uint) (*entity.GuaranteeViolation, error)
	CreateGuaranteeViolation(db database.Database, violation *entity.GuaranteeViolation) error
	UpdateGuaranteeViolation(db database.Database, violation *entity.GuaranteeViolation) error
	DeleteGuaranteeViolation(db database.Database, violation *entity.GuaranteeViolation) error
}
