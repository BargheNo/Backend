package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type CorporationStaff struct {
	database.Model
	UserID        uint   `gorm:"not null;index"`
	CorporationID uint   `gorm:"not null;index"`
	UserType      string `gorm:"type:varchar(50)"`
}
