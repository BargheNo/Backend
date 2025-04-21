package service

import (
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
)

type CorporationService interface {
	DoesCorporationExist(corporationID uint)
	ISCorporationApproved(corporationID uint) bool
	GetCorporationCredentials(corporationID uint) corporationdto.CorporationCredentialResponse
	CheckApplicantAccess(corporationID, applicantID uint)
	Register(registerInfo corporationdto.RegisterRequest) corporationdto.CorporationCredentialResponse
	UpdateRegister(updateRegisterInfo corporationdto.UpdateRegisterRequest)
	AddCertificateFiles(requestInfo corporationdto.AddCertificatesRequest)
	AddContactInfo(contactInfo corporationdto.AddContactInformationRequest)
	UpdateContactInfo(contactInfo corporationdto.AddContactInformationRequest)
	AddAddress(addressInfo corporationdto.AddCorporationAddressRequest)
	UpdateAddress(addressInfo corporationdto.AddCorporationAddressRequest)
	DeleteAddress(addressInfo corporationdto.DeleteAddressRequest)
	GetCorporationDetails(requestInfo corporationdto.CorporationDetailsRequest) corporationdto.CorporationPrivateInfoResponse
	GetContactInfo(corporationID uint) []corporationdto.ContactInformationResponse
	GetContactTypes() []corporationdto.ContactTypeResponse
}
