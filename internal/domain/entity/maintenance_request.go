package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type MaintenanceRequest struct {
	database.Model
	PanelID              uint                          `gorm:"index"`
	Panel                Panel                         `gorm:"foreignKey:PanelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CorporationID        uint                          `gorm:"index"`
	Corporation          Corporation                   `gorm:"foreignKey:CorporationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Status               enum.MaintenanceRequestStatus `gorm:"not null"`
	Subject              string                        `gorm:"type:varchar(50);not null"`
	Description          string                        `gorm:"type:text"`
	UrgencyLevel         enum.UrgencyLevel             `gorm:"not null"`
	IsGuaranteeRequested bool                          `gorm:"not null"`
}
