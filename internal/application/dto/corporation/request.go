package corporationdto

type RegisterRequest struct {
	Name     string
	CIN      string
	Password string
}

type LoginRequest struct {
	CIN      string
	Password string
}

type ContactInfoRequest struct {
	Phone     string
	Email     string
	Eitaa     string
	Bale      string
	Website   string
	WhatsApp  string
	Instagram string
	Linkedin  string
	Telegram  string
}

type AddressRequest struct {
	CorporationID  uint
	Province       string
	City           string
	StreetAddress  string
	PostalCode     string
	BuildingNumber string
	Unit           uint
}
