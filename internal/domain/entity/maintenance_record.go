package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type MaintenanceRecord struct {
	database.Model
	OperatorID           uint                `gorm:"not null"`
	Operator             User                `gorm:"foreignKey:OperatorID;references:ID"`
	RequestID            uint                `gorm:"not null;index"`
	Request              MaintenanceRequest  `gorm:"foreignKey:RequestID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IsUserApproved       bool                `gorm:"not null;default:false"`
	Title                string              `gorm:"not null"`
	Details              string              `gorm:"not null"`
	GuaranteeViolationID *uint               `gorm:"index"`
	GuaranteeViolation   *GuaranteeViolation `gorm:"foreignKey:GuaranteeViolationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
