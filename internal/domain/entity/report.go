package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Report struct {
	database.Model
	Description    string            `gorm:"type:text;not null"`
	ObjectID       uint              `gorm:"not null;index"`
	ObjectType     string            `gorm:"type:varchar(50);not null"`
	ReportedByID   uint              `gorm:"not null;index"`
	ReportedByType string            `gorm:"type:varchar(50);not null"`
	Status         enum.ReportStatus `gorm:"type:varchar(50);not null"`
}
