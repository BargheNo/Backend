package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)
type Corporation struct {
	database.Model
	Name    string
	CIN     uint
	Address []Address `gorm:"polymorphic:Owner;"`
}
