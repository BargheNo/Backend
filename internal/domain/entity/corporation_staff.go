package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type CorporationStaff struct {
	database.Model
	UserID        uint             `gorm:"not null;index"`
	User          User             `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CorporationID uint             `gorm:"not null;index"`
	Corporation   Corporation      `gorm:"foreignKey:CorporationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status        enum.StaffStatus `gorm:"default:1;not null"`
	Roles         []Role           `gorm:"many2many:staff_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
