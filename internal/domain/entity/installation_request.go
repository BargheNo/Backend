package entity

import "github.com/BargheNo/Backend/internal/infrastructure/database"

type InstallationRequest struct {
	database.Model  
	UserID         	uint  
	Area           	float64 
	PowerRequested  float64 
	MaxCost        	float64 
	Deadline       	string
	BuildingType   	string  
	Status			string
	Address        	string
}