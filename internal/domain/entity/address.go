package entity

import (
	"fmt"

	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type Address struct {
	database.Model
	Province       string
	City           string
	StreetAddress  string
	PostalCode     string
	BuildingNumber string
	Unit           uint
	OwnerID        uint
	OwnerType      string
}

func (a Address) String() string {
	return fmt.Sprintf("Province: %s, City: %s, "+
		"StreetAddress: %s, PostalCode: %s, "+
		"BuildingNumber: %s, Unit: %d",
		a.Province, a.City,
		a.StreetAddress, a.PostalCode,
		a.BuildingNumber, a.Unit)
}
