package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type PanelHistory struct {
	database.Model
	DatalogSerial string  `gorm:"type:varchar(50)"`
	PVSerial      string  `gorm:"type:varchar(50)"`
	Date          string  `gorm:"type:date;not null;index"`
	EnergyToday   float64 `gorm:"not null"`
	EnergyTotal   float64 `gorm:"not null"`
	PanelID       uint    `gorm:"not null;index"`
	Panel         Panel   `gorm:"foreignKey:PanelID"`
}
