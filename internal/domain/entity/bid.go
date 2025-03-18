package entity

import (
	"time"

	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Bid struct {
	database.Model
	CorporationID       uint `gorm:"not null"`
	RequestID           uint `gorm:"not null"`
	MinCost             float64
	MaxCost             float64
	MinDeadline         time.Time
	MaxDeadline         time.Time
	Description         string
	Status              string
	InstallationTime    string
}
