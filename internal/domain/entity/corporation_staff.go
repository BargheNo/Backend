package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type CorporationStaff struct {
	database.Model
	StaffID       uint   `gorm:"not null;index"`
	CorporationID uint   `gorm:"not null;index"`
	Roles         []Role `gorm:"many2many:staff_roles;constraint:OnDelete:CASCADE;"`
}
