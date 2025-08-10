package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type PaymentRepository interface {
	FindPaymentTerms(db database.Database, payTermID uint) (*entity.PaymentTerm, error)
	FindPaymentTermInstallmentPlan(db database.Database, payTermID uint) (*entity.InstallmentPlan, error)
	CreatePaymentTerms(db database.Database, paymentTerms *entity.PaymentTerm) error
	CreateInstallmentPlan(db database.Database, plan *entity.InstallmentPlan) error
	UpdatePaymentTerms(db database.Database, paymentTerms *entity.PaymentTerm) error
	UpdateInstallmentPlan(db database.Database, plan *entity.InstallmentPlan) error
}
