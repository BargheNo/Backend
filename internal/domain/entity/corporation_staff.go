package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type CorporationStaff struct {
	database.Model
	StaffID       uint           `gorm:"not null;index"`
	CorporationID uint           `gorm:"not null;index"`
	StaffType     enum.StaffType `gorm:"not null;index"`
}
