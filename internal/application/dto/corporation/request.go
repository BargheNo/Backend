package corporationdto

import (
	"mime/multipart"

	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
)

type Signatory struct {
	Name               string
	NationalCardNumber string
	Position           string
}

type RegisterRequest struct {
	ApplicantID        uint
	Name               string
	NationalID         string
	RegistrationNumber string
	IBAN               string
	Signatories        []Signatory
}

type UpdateRegisterRequest struct {
	ApplicantID        uint
	CorporationID      uint
	Name               *string
	NationalID         *string
	RegistrationNumber *string
	IBAN               *string
	Signatories        []Signatory
}

type AddCertificatesRequest struct {
	CorporationID          uint
	ApplicantID            uint
	VATTaxpayerCertificate *multipart.FileHeader
	OfficialNewspaperAD    *multipart.FileHeader
}

type ContactInformation struct {
	ContactTypeID uint
	ContactValue  string
}

type AddContactInformationRequest struct {
	ApplicantID        uint
	CorporationID      uint
	ContactInformation []ContactInformation
}

type AddCorporationAddressRequest struct {
	ApplicantID   uint
	CorporationID uint
	Addresses     []addressdto.CreateAddressRequest
}

type DeleteAddressRequest struct {
	UserID        uint
	CorporationID uint
	AddressID     uint
}

type CorporationListRequest struct {
	UserID uint
	Offset int
	Limit  int
}
