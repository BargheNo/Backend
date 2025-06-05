package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Corporation struct {
	database.Model
	Name                   string                 `gorm:"type:varchar(100);unique;not null"`
	Logo                   string                 `gorm:"type:text"`
	RegistrationNumber     string                 `gorm:"type:varchar(50);unique;not null"`
	NationalID             string                 `gorm:"type:varchar(50);unique;not null"`
	VATTaxpayerCertificate string                 `gorm:"type:varchar(255)"`
	OfficialNewspaperAD    string                 `gorm:"type:varchar(255)"`
	IBAN                   string                 `gorm:"type:varchar(34)"`
	Status                 enum.CorporationStatus `gorm:"index"`
	Signatories            []Signatory            `gorm:"foreignKey:CorporationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ContactInformation     []ContactInformation   `gorm:"foreignKey:CorporationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Addresses              []Address              `gorm:"polymorphic:Owner;polymorphicValue:corporations"`
	Bids                   []Bid                  `gorm:"foreignKey:CorporationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Guarantees             []Guarantee            `gorm:"foreignKey:CorporationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
