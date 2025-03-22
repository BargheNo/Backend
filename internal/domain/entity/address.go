package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Address struct {
	database.Model
	ProvinceID    uint
	Province      Province `gorm:"foreignKey:ProvinceID"`
	CityID        uint
	City          City `gorm:"foreignKey:CityID"`
	StreetAddress string
	PostalCode    string
	HouseNumber   string
	Unit          uint
	OwnerID       uint
	OwnerType     string
}
