package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Address struct {
	database.Model
	ProvinceID    uint     `gorm:"not null;index"`
	Province      Province `gorm:"foreignKey:ProvinceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	CityID        uint     `gorm:"not null;index"`
	City          City     `gorm:"foreignKey:CityID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	StreetAddress string   `gorm:"type:text;not null"`
	PostalCode    string   `gorm:"type:varchar(20);not null"`
	HouseNumber   string   `gorm:"type:varchar(50);not null"`
	Unit          uint     `gorm:"default:0"`
	OwnerID       uint     `gorm:"not null;index"`
	OwnerType     string   `gorm:"type:varchar(50);not null"`
}
