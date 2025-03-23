package entity

import (
	"time"

	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Bid struct {
	database.Model
	CorporationID    uint `gorm:"not null"`
	RequestID        uint `gorm:"not null"`
	Request          InstallationRequest
	Cost             uint
	Description      string
	Status           enum.BidStatus
	InstallationDate time.Time
}
