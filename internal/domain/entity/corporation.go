package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Corporation struct {
	database.Model
	Name                 string
	CIN                  string
	Status               string
	Password             string
	ContactInformation   ContactInformation
	Addresses            []Address `gorm:"polymorphic:Owner;"`
	Bids                 []Bid	`gorm:"foreignKey:CorporationID"`
}
