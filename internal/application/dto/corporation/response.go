package corporationdto

import (
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
)

type CorporationDetailsResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CorporationLoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Name         string `json:"name"`
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
	ID          uint                         `json:"id"`
	Name        string                       `json:"name"`
	ContactInfo []ContactInformationResponse `json:"contactInfo"`
	Addresses   []addressdto.AddressResponse `json:"addresses"`
}

type ContactInformationResponse struct {
	ContactTypeID uint
	ContactValue  string
}
