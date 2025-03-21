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

type ListOwnerRequestsRequest struct {
	OwnerID uint
	Offset  int
	Limit   int
}
