package entity

import (
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type CorporationReview struct {
	database.Model
	CorporationID uint              `gorm:"not null;index"`
	Corporation   Corporation       `gorm:"foreignKey:CorporationID"`
	ReviewerID    uint              `gorm:"not null;index"`
	Reviewer      User              `gorm:"foreignKey:ReviewerID"`
	Action        enum.ReviewAction `gorm:"not null"`
	Reason        *string           `gorm:"type:text"`
	Notes         *string           `gorm:"type:text"`
}
