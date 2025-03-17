package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type InstallationRequest struct {
	database.Model
	Name         string
	Status       string
	OwnerID      uint
	Owner        User
	Address      []Address `gorm:"polymorphic:Owner;"`
	Area         uint
	PowerRequest uint
	MaxCost      uint
	BuildingType string
}
