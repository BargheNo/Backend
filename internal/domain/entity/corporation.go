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
	ContactInformationID uint
	Address              []Address          `gorm:"polymorphic:Owner;"`
	ContactInformation   ContactInformation `gorm:"foreignKey:ContactInformationID"`
}
