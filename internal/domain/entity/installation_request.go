package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type InstallationRequest struct {
	database.Model
	Name         string                         `gorm:"type:varchar(50);not null"`
	Status       enum.InstallationRequestStatus `gorm:"not null"`
	PowerRequest uint                           `gorm:"not null"`
	Description  string                         `gorm:"type:text"`
	BuildingType enum.BuildingType              `gorm:"not null"`
	OwnerID      uint                           `gorm:"index"`
	Owner        User                           `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Address      Address                        `gorm:"polymorphic:Owner;polymorphicValue:installation_requests"`
	Area         uint
	MaxCost      float64
}
