package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

// you can add national card photo path to Signatory too
type Signatory struct {
	database.Model
	CorporationID      uint   `gorm:"index"`
	Name               string `gorm:"type:varchar(50);not null"`
	NationalCardNumber string `gorm:"type:varchar(50);not null"`
	Position           string `gorm:"type:varchar(100)"`
}
