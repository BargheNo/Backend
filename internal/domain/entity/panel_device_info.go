package entity

import (
	"time"

	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type PanelDeviceInfo struct {
	database.Model

	DatalogSerial string `gorm:"type:varchar(50);not null"`
	PVSerial      string `gorm:"type:varchar(50);not null"`
	DeviceModel   string `gorm:"type:varchar(50)"`
	Firmware      string `gorm:"type:varchar(20)"`
	Hardware      string `gorm:"type:varchar(20)"`
	LoggerFW      string `gorm:"type:varchar(20)"`
	Timezone      string `gorm:"type:varchar(20)"`

	LastUpdated time.Time `gorm:"not null"`

	PanelID uint  `gorm:"not null;index"`
	Panel   Panel `gorm:"foreignKey:PanelID"`
}
