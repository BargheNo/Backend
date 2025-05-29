package installationdto

import (
	"time"

	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type EnumStatusResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AnonymousRequestsResponse struct {
	ID           uint                       `json:"id"`
	Name         string                     `json:"name"`
	CreatedTime  time.Time                  `json:"createdTime"`
	Status       string                     `json:"status"`
	PowerRequest uint                       `json:"powerRequest"`
	MaxCost      float64                    `json:"maxCost"`
	BuildingType string                     `json:"buildingType"`
	Address      addressdto.AddressResponse `json:"address"`
}

type PublicRequestDetailsResponse struct {
	ID           uint                       `json:"id"`
	Name         string                     `json:"name"`
	Status       string                     `json:"status"`
	PowerRequest uint                       `json:"powerRequest"`
	Description  string                     `json:"description"`
	BuildingType string                     `json:"buildingType"`
	Area         uint                       `json:"area"`
	MaxCost      float64                    `json:"maxCost"`
	Customer     userdto.CredentialResponse `json:"customer"`
	Address      addressdto.AddressResponse `json:"address"`
}

type CorporationPanelResponse struct {
	ID                   uint                           `json:"id"`
	Name                 string                         `json:"name"`
	Status               string                         `json:"status"`
	BuildingType         string                         `json:"buildingType"`
	Area                 uint                           `json:"area"`
	Power                uint                           `json:"power"`
	Tilt                 uint                           `json:"tilt"`
	Azimuth              uint                           `json:"azimuth"`
	TotalNumberOfModules uint                           `json:"totalNumberOfModules"`
	Operator             userdto.CredentialResponse     `json:"operator"`
	Customer             userdto.CredentialResponse     `json:"customer"`
	Address              addressdto.AddressResponse     `json:"address"`
	Guarantee            guaranteedto.GuaranteeResponse `json:"guarantee"`
}

type PanelResponse struct {
	ID                   uint                                         `json:"id"`
	Name                 string                                       `json:"name"`
	Status               string                                       `json:"status"`
	BuildingType         string                                       `json:"buildingType"`
	Area                 uint                                         `json:"area"`
	Power                uint                                         `json:"power"`
	Tilt                 uint                                         `json:"tilt"`
	Azimuth              uint                                         `json:"azimuth"`
	TotalNumberOfModules uint                                         `json:"totalNumberOfModules"`
	Operator             userdto.CredentialResponse                   `json:"operator"`
	Customer             userdto.CredentialResponse                   `json:"customer"`
	Corporation          corporationdto.CorporationCredentialResponse `json:"corporation"`
	Address              addressdto.AddressResponse                   `json:"address"`
	Guarantee            guaranteedto.GuaranteeResponse               `json:"guarantee"`
}

type CustomerPanelResponse struct {
	ID                   uint                                         `json:"id"`
	Name                 string                                       `json:"name"`
	Status               string                                       `json:"status"`
	BuildingType         string                                       `json:"buildingType"`
	Area                 uint                                         `json:"area"`
	Power                uint                                         `json:"power"`
	Tilt                 uint                                         `json:"tilt"`
	Azimuth              uint                                         `json:"azimuth"`
	TotalNumberOfModules uint                                         `json:"totalNumberOfModules"`
	Corporation          corporationdto.CorporationCredentialResponse `json:"corporation"`
	Address              addressdto.AddressResponse                   `json:"address"`
	Guarantee            guaranteedto.GuaranteeResponse               `json:"guarantee"`
}
