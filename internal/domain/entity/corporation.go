package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enums"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Corporation struct {
	database.Model
	Name               string
	CIN                string
	Status             enums.CorporationStatus
	Password           string
	ContactInformation ContactInformation
	Addresses          []Address `gorm:"polymorphic:Owner;"`
	Bids               []Bid     `gorm:"foreignKey:CorporationID"`
}
