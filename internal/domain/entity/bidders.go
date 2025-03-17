package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type Bidders struct {
	database.Model	
	CorporationID		uint
	RequestID			uint
	RequestType			string
	MinCost				float64
	MaxCost				float64
	MinDeadline			string
	MaxDeadline			string
	Description			string
	Status				string
	InstallationTime	string
}