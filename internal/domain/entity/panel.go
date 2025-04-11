package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type Panel struct {
	database.Model
	Name                 string      `gorm:"type:varchar(50);not null"`
	BuildingType         string      `gorm:"type:varchar(50);not null"`
	Area                 uint        `gorm:"not null"`
	Power                uint        `gorm:"not null"`
	Tilt                 uint        `gorm:"not null"`
	Azimuth              uint        `gorm:"not null"`
	TotalNumberOfModules uint        `gorm:"not null"`
	CorporationID        uint        `gorm:"not null;index"`
	Corporation          Corporation `gorm:"foreignKey:CorporationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	OperatorID           uint        `gorm:"not null;index"`
	Operator             User        `gorm:"foreignKey:OperatorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CustomerID           uint        `gorm:"not null;index"`
	Customer             User        `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Address              Address     `gorm:"polymorphic:Owner;polymorphicValue:panels"`
}
