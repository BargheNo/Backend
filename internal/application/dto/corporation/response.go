package corporationdto

import (
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
)

type CorporationCredentialResponse struct {
	ID          uint                         `json:"id"`
	Name        string                       `json:"name"`
	Logo        string                       `json:"logo"`
	ContactInfo []ContactInformationResponse `json:"contactInfo"`
	Addresses   []addressdto.AddressResponse `json:"addresses"`
}

type CorporationPrivateInfoResponse struct {
	ID                     uint                         `json:"id"`
	Name                   string                       `json:"name"`
	RegistrationNumber     string                       `json:"registrationNumber"`
	NationalID             string                       `json:"nationalID"`
	IBAN                   string                       `json:"iban"`
	Logo                   string                       `json:"logo"`
	VATTaxpayerCertificate string                       `json:"vatTaxpayerCertificate"`
	OfficialNewspaperAD    string                       `json:"officialNewspaperAD"`
	Signatories            []SignatoryResponse          `json:"signatories"`
	ContactInfo            []ContactInformationResponse `json:"contactInfo"`
	Addresses              []addressdto.AddressResponse `json:"addresses"`
}

type SignatoryResponse struct {
	ID                 uint   `json:"id"`
	Name               string `json:"name"`
	NationalCardNumber string `json:"nationalCardNumber"`
	Position           string `json:"position"`
}

type ContactInformationResponse struct {
	ID          uint
	ContactType ContactTypeResponse `json:"contactType"`
	Value       string              `json:"value"`
}

type ContactTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
