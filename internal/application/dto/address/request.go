package addressdto

type CreateAddressRequest struct {
	ProvinceID    uint
	CityID        uint
	StreetAddress string
	PostalCode    string
	HouseNumber   string
	Unit          uint
	OwnerID       uint
	OwnerType     string
}

type GetOwnerAddressesRequest struct {
	OwnerID   uint
	OwnerType string
}

type GetProvinceCitiesRequest struct {
	ProvinceID uint
}
