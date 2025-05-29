package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type PaymentTerm struct {
	database.Model
	PaymentMethod   enum.PaymentMethod `gorm:"not null"`
	InstallmentPlan *InstallmentPlan   `gorm:"foreignKey:PaymentTermsID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type InstallmentPlan struct {
	database.Model
	PaymentTermsID    uint   `gorm:"not null;index"`
	NumberOfMonths    uint   `gorm:"not null"`
	DownPaymentAmount uint   `gorm:"not null"`
	MonthlyAmount     uint   `gorm:"not null"`
	Notes             string `gorm:"type:text"`
	// DownPaymentDate   string `gorm:"type:varchar(50)"`
	// DueDay            uint   `gorm:"not null"`
}
