package entity

import (
	"time"

	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type PanelEvent struct {
	database.Model
	DatalogSerial string       `gorm:"type:varchar(50)"`
	PVSerial      string       `gorm:"type:varchar(50)"`
	EventCode     string       `gorm:"type:varchar(10);not null"`
	Description   string       `gorm:"type:varchar(255);not null"`
	Severity      string       `gorm:"type:varchar(20);not null;index"`
	Timestamp     time.Time    `gorm:"not null;index"`
	PanelID       uint         `gorm:"not null;index"`
	Panel         entity.Panel `gorm:"foreignKey:PanelID"`
}
