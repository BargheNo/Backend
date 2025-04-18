package corporationdto

import (
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
)

type CorporationDetailsResponse struct {
	ID          uint                         `json:"id"`
	Name        string                       `json:"name"`
	Logo        string                       `json:"logo"`
	ContactInfo []ContactInformationResponse `json:"contactInfo"`
	Addresses   []addressdto.AddressResponse `json:"addresses"`
}

type CorporationLoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Name         string `json:"name"`
}

type ContactInformationResponse struct {
	ContactTypeID uint
	ContactValue  string
}
