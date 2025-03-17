package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type Bidders struct {
	database.Model	
	CorporationID		Corporation `gorm:"foreignKey:CorporationID"`
	RequestID			InstallationRequest `gorm:"foreignKey:RequestID"`
	RequestType			string
	MinCost				float64
	MaxCost				float64
	MinDeadline			string
	MaxDeadline			string
	Description			string
	Status				string
	InstallationTime	string
}