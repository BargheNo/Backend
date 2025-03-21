package installationdto

type NewInstallationRequest struct {
	OwnerID      uint
	Name         string
	Area         uint
	Power        uint
	MaxCost      float64
	BuildingType string
	Description  string
	AddressID    uint
}

type ListOwnerRequestsRequest struct {
	OwnerID uint
	Offset  int
	Limit   int
}
