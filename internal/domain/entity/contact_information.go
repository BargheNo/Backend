package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type ContactType struct {
	database.Model
	TypeName string `gorm:"type:varchar(255);not null;unique"`
}

type ContactInformation struct {
	database.Model
	CorporationID uint   `gorm:"not null"`
	ContactTypeID uint   `gorm:"not null"`
	ContactValue  string `gorm:"type:varchar(255);not null"`
	IsPrimary     bool   `gorm:"default:false"`
}
