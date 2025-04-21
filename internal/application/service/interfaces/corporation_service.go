package service

import (
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
)

type CorporationService interface {
	DoesCorporationExist(corporationID uint)
	ISCorporationApproved(corporationID uint) bool
	GetCorporationCredentials(corporationID uint) corporationdto.CorporationDetailsResponse
	CheckApplicantAccess(corporationID, applicantID uint)
	Register(registerInfo corporationdto.RegisterRequest) corporationdto.CorporationDetailsResponse
	UpdateRegister(updateRegisterInfo corporationdto.UpdateRegisterRequest)
	AddCertificateFiles(requestInfo corporationdto.AddCertificatesRequest)
	AddContactInfo(contactInfo corporationdto.AddContactInformationRequest)
	UpdateContactInfo(contactInfo corporationdto.AddContactInformationRequest)
	AddAddress(addressInfo corporationdto.AddCorporationAddressRequest)
	DeleteAddress(addressInfo corporationdto.DeleteAddressRequest)
	GetCorporations(requestInfo corporationdto.CorporationListRequest) []corporationdto.CorporationDetailsResponse
	GetContactInfo(corporationID uint) []corporationdto.ContactInformationResponse
}
