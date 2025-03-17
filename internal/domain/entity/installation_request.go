package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type InstallationRequest struct {
	database.Model
	UserID         uint
	User           User `gorm:"foreignKey:UserID"`
	Area           float64
	PowerRequested float64
	MaxCost        float64
	Deadline       string
	BuildingType   string
	Status         string
	Address        Address `gorm:"polymorphic:Owner;"`
	Bids		  []Bid `gorm:"foreignKey:RequestID"`
}
