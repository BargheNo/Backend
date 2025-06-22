package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type PaymentRepository struct{}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{}
}

func (repo *PaymentRepository) FindPaymentTerms(db database.Database, payTermID uint) (*entity.PaymentTerm, error) {
	var payTerm *entity.PaymentTerm
	result := db.GetDB().First(&payTerm, payTermID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return payTerm, nil
}

func (repo *PaymentRepository) FindPaymentTermInstallmentPlan(db database.Database, payTermID uint) (*entity.InstallmentPlan, error) {
	var plan *entity.InstallmentPlan
	result := db.GetDB().Where("payment_terms_id = ?", payTermID).First(&plan)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return plan, nil
}

func (repo *PaymentRepository) CreatePaymentTerms(db database.Database, paymentTerms *entity.PaymentTerm) error {
	return db.GetDB().Create(&paymentTerms).Error
}

func (repo *PaymentRepository) CreateInstallmentPlan(db database.Database, plan *entity.InstallmentPlan) error {
	return db.GetDB().Create(&plan).Error
}

func (repo *PaymentRepository) UpdatePaymentTerms(db database.Database, paymentTerms *entity.PaymentTerm) error {
	return db.GetDB().Save(&paymentTerms).Error
}

func (repo *PaymentRepository) UpdateInstallmentPlan(db database.Database, plan *entity.InstallmentPlan) error {
	return db.GetDB().Save(&plan).Error
}
