package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Guarantee struct {
	database.Model
	CorporationID  uint                 `gorm:"not null;index"`
	Corporation    Corporation          `gorm:"foreignKey:CorporationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name           string               `gorm:"type:varchar(100);not null"`
	Status         enum.GuaranteeStatus `gorm:"not null"`
	GuaranteeType  enum.GuaranteeType   `gorm:"not null"`
	DurationMonths uint                 `gorm:"not null"`
	Description    string               `gorm:"type:text"`
	Terms          []GuaranteeTerm      `gorm:"foreignKey:GuaranteeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type GuaranteeTerm struct {
	database.Model
	GuaranteeID uint   `gorm:"not null;index"`
	Title       string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:text;not null"`
	Limitations string `gorm:"type:text"`
}
