package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type GuaranteeViolation struct {
	database.Model
	PanelID      uint   `gorm:"index"`
	Panel        Panel  `gorm:"foreignKey:PanelID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ViolatedByID uint   `gorm:"not null"`
	ViolatedBy   User   `gorm:"foreignKey:ViolatedByID"`
	Reason       string `gorm:"type:text;not null"`
	Details      string `gorm:"type:text"`
}
