package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Corporation struct {
	database.Model
	Name                   string                 `gorm:"type:varchar(100);unique;not null"`
	RegistrationNumber     string                 `gorm:"type:varchar(50);unique;not null"`
	NationalID             string                 `gorm:"type:varchar(50);unique;not null"`
	VATTaxpayerCertificate *string                `gorm:"type:varchar(255)"`
	OfficialNewspaperAD    *string                `gorm:"type:varchar(255)"`
	CompanyIBAN            *string                `gorm:"type:varchar(34)"`
	Signatories            []Signatory            `gorm:"foreignKey:CorporationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ContactInformation     []ContactInformation   `gorm:"foreignKey:CorporationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Addresses              []Address              `gorm:"polymorphic:Owner;polymorphicValue:corporation"`
	Status                 enum.CorporationStatus `gorm:"type:varchar(20);index"`
	Bids                   []Bid                  `gorm:"foreignKey:CorporationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// you can add national card photo path to Signatory too
type Signatory struct {
	database.Model
	CorporationID      uint   `gorm:"index"`
	Name               string `gorm:"type:varchar(100);not null"`
	NationalCardNumber string `gorm:"type:varchar(50);not null"`
	Position           string `gorm:"type:varchar(100)"`
}
