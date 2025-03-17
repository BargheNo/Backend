package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type Bid struct {
	database.Model
	CorporationID       uint
	Corporation         Corporation `gorm:"foreignKey:CorporationID"`
	RequestID           uint
	InstallationRequest InstallationRequest `gorm:"foreignKey:RequestID"`
	MinCost             float64
	MaxCost             float64
	MinDeadline         string
	MaxDeadline         string
	Description         string
	Status              string
	InstallationTime    string
}
