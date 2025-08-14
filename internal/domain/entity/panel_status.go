package entity

import (
	"time"

	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type PanelStatus struct {
	database.Model
	DatalogSerial string    `gorm:"type:varchar(50)"`
	PVSerial      string    `gorm:"type:varchar(50)"`
	PVStatus      int       `gorm:"not null"`
	PVPowerIn     float64   `gorm:"not null"`
	PV1Voltage    float64   `gorm:"not null"`
	PV1Current    float64   `gorm:"not null"`
	PV2Voltage    float64   `gorm:"not null"`
	PV2Current    float64   `gorm:"not null"`
	PVPowerOut    float64   `gorm:"not null"`
	ACFreq        float64   `gorm:"not null"`
	ACVoltage     float64   `gorm:"not null"`
	ACOutputPower float64   `gorm:"not null"`
	Temperature   float64   `gorm:"not null"`
	BatVoltage    float64   `gorm:"not null"`
	BatCurrent    float64   `gorm:"not null"`
	BatPower      float64   `gorm:"not null"`
	GridExport    float64   `gorm:"not null"`
	GridImport    float64   `gorm:"not null"`
	EnergyToday   float64   `gorm:"not null"`
	EnergyTotal   float64   `gorm:"not null"`
	Timestamp     time.Time `gorm:"not null;index"`
	PanelID       uint      `gorm:"not null;index"`
	Panel         Panel     `gorm:"foreignKey:PanelID"`
}
