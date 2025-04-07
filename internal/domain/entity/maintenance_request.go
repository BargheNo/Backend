package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type MaintenanceRequest struct {
	database.Model
	Subject       string                        `gorm:"type:varchar(50);not null"`
	Description   string                        `gorm:"type:text"`
	PanelID       uint                          `gorm:"index"`
	Panel         Panel                         `gorm:"foreignKey:PanelID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	OwnerID       uint                          `gorm:"index"`
	Owner         User                          `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CorporationID uint                          `gorm:"index"`
	Corporation   Corporation                   `gorm:"foreignKey:CorporationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UrgencyLevel  enum.UrgencyLevel             `gorm:"not null"`
	Status        enum.MaintenanceRequestStatus `gorm:"not null"`
}
