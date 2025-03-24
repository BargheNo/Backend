package installationdto

import (
	"time"

	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type OwnerRequestsResponse struct {
	ID           uint                       `json:"id"`
	Name         string                     `json:"name"`
	CreatedTime  time.Time                  `json:"createdTime"`
	Status       string                     `json:"status"`
	PowerRequest uint                       `json:"powerRequest"`
	MaxCost      float64                    `json:"maxCost"`
	BuildingType string                     `json:"buildingType"`
	Address      addressdto.AddressResponse `json:"address"`
}

type RequestDetailsResponse struct {
	ID           uint                       `json:"id"`
	Name         string                     `json:"name"`
	CreatedTime  time.Time                  `json:"createdTime"`
	Status       string                     `json:"status"`
	PowerRequest uint                       `json:"powerRequest"`
	MaxCost      float64                    `json:"maxCost"`
	BuildingType string                     `json:"buildingType"`
	Address      addressdto.AddressResponse `json:"address"`
	Customer     userdto.CredentialResponse `json:"customer"`
}
