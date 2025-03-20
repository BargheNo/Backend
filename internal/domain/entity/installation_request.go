package entity

import (
	"time"

	"github.com/BargheNo/Backend/internal/domain/enums"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type InstallationRequest struct {
	database.Model
	UserID         uint `gorm:"not null"`
	Area           float64
	PowerRequested float64
	MaxCost        float64
	Deadline       time.Time
	BuildingType   string
	Status         enums.InstallationRequestStatus
	Address        Address `gorm:"polymorphic:Owner;"`
	Bids           []Bid   `gorm:"foreignKey:RequestID"`
}
