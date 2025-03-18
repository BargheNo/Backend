package addressdto

type CreateAddressRequest struct {
	Province      string
	City          string
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
