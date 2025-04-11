package entity

import (
	"time"

	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type MaintenanceRecord struct {
	database.Model
	PanelID       uint        `json:"panel_id" gorm:"column:panel_id;not null"`
	Panel         Panel       `json:"panel" gorm:"foreignKey:PanelID;references:ID"`
	CorporationID uint        `json:"corporation_id" gorm:"column:corporation_id;not null"`
	Corporation   Corporation `json:"corporation" gorm:"foreignKey:CorporationID;references:ID"`
	CustomerID    uint        `json:"customer_id" gorm:"column:customer_id;not null"`
	Customer      User        `json:"customer" gorm:"foreignKey:CustomerID;references:ID"`
	OperatorID    uint        `json:"operator_id" gorm:"column:operator_id;not null"`
	Operator      User        `json:"operator" gorm:"foreignKey:OperatorID;references:ID"`
	Title         string      `json:"title" gorm:"column:title;not null"`
	Details       string      `json:"details" gorm:"column:details;not null"`
	Date          time.Time   `json:"date" gorm:"column:date;not null"`
}
