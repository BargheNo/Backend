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

type SetBidRequest struct {
	InstallationRequestID uint
	CorporationID         uint
	MinCost               float64
	MaxCost               float64
	MinDeadline           string
	MaxDeadline           string
	Description           string
	InstallationTime      string
}

type CancelBidRequest struct {
	BidID                 uint
	InstallationRequestID uint
	CorporationID         uint
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
