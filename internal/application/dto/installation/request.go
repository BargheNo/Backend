package installationdto

import addressdto "github.com/BargheNo/Backend/internal/application/dto/address"

type NewInstallationRequest struct {
	OwnerID      uint
	Name         string
	Area         uint
	Power        uint
	MaxCost      float64
	BuildingType string
	Description  string
	Address      addressdto.CreateAddressRequest
}

type InstallationListRequest struct {
	OwnerID uint
	Offset  int
	Limit   int
}

type PanelListRequest struct {
	CorporationID uint
	Offset        int
	Limit         int
}

type GetOwnerRequest struct {
	OwnerID   uint
	RequestID uint
}

type AddPanelRequest struct {
	CorporationID        uint
	OperatorID           uint
	PanelName            string
	CustomerPhone        string
	Power                uint
	Area                 uint
	BuildingType         string
	Tilt                 uint
	Azimuth              uint
	TotalNumberOfModules uint
	Address              addressdto.CreateAddressRequest
}
