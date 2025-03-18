package addressdto

type AddressResponse struct {
	ID            uint   `json:"ID"`
	Province      string `json:"province"`
	City          string `json:"city"`
	StreetAddress string `json:"streetAddress"`
	PostalCode    string `json:"postalCode"`
	HouseNumber   string `json:"houseNumber"`
	Unit          uint   `json:"unit"`
}
