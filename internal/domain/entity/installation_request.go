package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type InstallationRequest struct {
	database.Model
	Name         string
	Status       enum.InstallationRequestStatus
	Area         uint
	PowerRequest uint
	MaxCost      float64
	BuildingType string
	OwnerID      uint
	Owner        User
	AddressID    uint
	Address      Address `gorm:"foreignKey:AddressID"`
}
