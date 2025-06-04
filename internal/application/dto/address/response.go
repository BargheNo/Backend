package addressdto

type AddressResponse struct {
	ID            uint   `json:"id"`
	Province      string `json:"province"`
	ProvinceID    uint   `json:"provinceID"`
	CityID        uint   `json:"cityID"`
	City          string `json:"city"`
	StreetAddress string `json:"streetAddress"`
	PostalCode    string `json:"postalCode"`
	HouseNumber   string `json:"houseNumber"`
	Unit          uint   `json:"unit"`
}

type ProvinceResponse struct {
	ID   uint   `json:"ID"`
	Name string `json:"name"`
}

type CityResponse struct {
	ID   uint   `json:"ID"`
	Name string `json:"name"`
}
