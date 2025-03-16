package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)
type Corporation struct {
	database.Model
	Name    	string
	CIN     	string
	Address 	[]Address `gorm:"polymorphic:Owner;"`
	Status		string
	Password 	string
}
