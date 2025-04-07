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

type CorporationPanelListResponse struct {
	ID                   uint                       `json:"id"`
	PanelName            string                     `json:"panelName"`
	CustomerName         string                     `json:"customerName"`
	CustomerPhone        string                     `json:"customerPhone"`
	Power                uint                       `json:"power"`
	Area                 uint                       `json:"area"`
	BuildingType         string                     `json:"buildingType"`
	Tilt                 uint                       `json:"tilt"`
	Azimuth              uint                       `json:"azimuth"`
	TotalNumberOfModules uint                       `json:"totalNumberOfModules"`
	Address              addressdto.AddressResponse `json:"address"`
	OperatorName         string                     `json:"operatorName"`
}

type CustomerPanelListResponse struct {
	ID                   uint                       `json:"id"`
	PanelName            string                     `json:"panelName"`
	CorporationName      string                     `json:"corporationName"`
	Power                uint                       `json:"power"`
	Area                 uint                       `json:"area"`
	BuildingType         string                     `json:"buildingType"`
	TotalNumberOfModules uint                       `json:"totalNumberOfModules"`
	Tilt                 uint                       `json:"tilt"`
	Azimuth              uint                       `json:"azimuth"`
	Address              addressdto.AddressResponse `json:"address"`
}
