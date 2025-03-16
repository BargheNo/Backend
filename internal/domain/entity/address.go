package entity

import (
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Address struct {
	database.Model
	Province      	string
	City          	string
	StreetAddress 	string
	PostalCode    	string
	BuildingNumber  string
	Unit          	uint
	OwnerID       	uint
	OwnerType     	string
}
