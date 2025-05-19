package installationdto

import (
	"time"

	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
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

type CorporationPanelResponse struct {
	ID                   uint                       `json:"id"`
	PanelName            string                     `json:"panelName"`
	Customer             userdto.CredentialResponse `json:"customer"`
	Operator             userdto.CredentialResponse `json:"operator"`
	Power                uint                       `json:"power"`
	Area                 uint                       `json:"area"`
	BuildingType         string                     `json:"buildingType"`
	Tilt                 uint                       `json:"tilt"`
	Azimuth              uint                       `json:"azimuth"`
	TotalNumberOfModules uint                       `json:"totalNumberOfModules"`
	Address              addressdto.AddressResponse `json:"address"`
}

type PanelResponse struct {
	ID                   uint                                         `json:"id"`
	Name                 string                                       `json:"name"`
	Customer             userdto.CredentialResponse                   `json:"customer"`
	Operator             userdto.CredentialResponse                   `json:"operator"`
	Corporation          corporationdto.CorporationCredentialResponse `json:"corporation"`
	Address              addressdto.AddressResponse                   `json:"address"`
	PanelName            string                                       `json:"panelName"`
	Power                uint                                         `json:"power"`
	Area                 uint                                         `json:"area"`
	BuildingType         string                                       `json:"buildingType"`
	Tilt                 uint                                         `json:"tilt"`
	Azimuth              uint                                         `json:"azimuth"`
	TotalNumberOfModules uint                                         `json:"totalNumberOfModules"`
}

type CustomerPanelResponse struct {
	ID                   uint                                         `json:"id"`
	PanelName            string                                       `json:"name"`
	Corporation          corporationdto.CorporationCredentialResponse `json:"corporation"`
	Power                uint                                         `json:"power"`
	Area                 uint                                         `json:"area"`
	BuildingType         string                                       `json:"buildingType"`
	TotalNumberOfModules uint                                         `json:"totalNumberOfModules"`
	Tilt                 uint                                         `json:"tilt"`
	Azimuth              uint                                         `json:"azimuth"`
	Address              addressdto.AddressResponse                   `json:"address"`
}
