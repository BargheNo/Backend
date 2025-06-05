package corporationdto

import (
	"mime/multipart"

	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	"github.com/BargheNo/Backend/internal/domain/enum"
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
	CorporationStatus  enum.CorporationStatus
	ContactInformation []ContactInformation
}

type DeleteContactInformationRequest struct {
	ApplicantID       uint
	ContactID         uint
	CorporationID     uint
	CorporationStatus enum.CorporationStatus
}

type AddCorporationAddressRequest struct {
	ApplicantID       uint
	CorporationID     uint
	CorporationStatus enum.CorporationStatus
	Addresses         []addressdto.CreateAddressRequest
}

type DeleteAddressRequest struct {
	UserID            uint
	CorporationID     uint
	CorporationStatus enum.CorporationStatus
	AddressID         uint
}

type CorporationDetailsRequest struct {
	UserID        uint
	CorporationID uint
	Status        enum.CorporationStatus
}

type ChangeLogoRequest struct {
	ApplicantID   uint
	CorporationID uint
	Logo          *multipart.FileHeader
}

type GetCorporationsByAdminRequest struct {
	Status uint
	Limit  int
	Offset int
}
