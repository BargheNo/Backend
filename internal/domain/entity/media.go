package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type Media struct {
	database.Model
	Path      string `gorm:"not null"`
	OwnerID   uint   `gorm:"not null;index"`
	OwnerType string `gorm:"type:varchar(50);not null"`
}
