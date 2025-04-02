package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type ContactType struct {
	database.Model
	Name string `gorm:"type:varchar(50);not null;unique"`
}

type ContactInformation struct {
	database.Model
	CorporationID uint   `gorm:"not null"`
	TypeID        uint   `gorm:"not null"`
	Value         string `gorm:"type:varchar(255);not null"`
}
