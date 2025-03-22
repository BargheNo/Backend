package corporationdto

type CorporationLoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Name         string `json:"name"`
}

type AddressResponse struct {
	ID             uint   `json:"id"`
	Province       string `json:"province"`
	City           string `json:"city"`
	StreetAddress  string `json:"streetAddress"`
	PostalCode     string `json:"postalCode"`
	BuildingNumber string `json:"buildingNumber"`
	Unit           uint   `json:"unit"`
}

type ContactInfoResponse struct {
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Eitaa     string `json:"eitaa"`
	Bale      string `json:"bale"`
	Website   string `json:"website"`
	WhatsApp  string `json:"whatsApp"`
	Instagram string `json:"instagram"`
	Linkedin  string `json:"linkedin"`
	Telegram  string `json:"telegram"`
}

type CorporationInfoResponse struct {
	ID          uint                `json:"id"`
	Name        string              `json:"name"`
	CIN         string              `json:"cin"`
	Status      string              `json:"status"`
	ContactInfo ContactInfoResponse `json:"contactInfo"`
	Addresses   []AddressResponse   `json:"addresses"`
}
