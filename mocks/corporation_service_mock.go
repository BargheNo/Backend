package mocks

import (
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	"github.com/stretchr/testify/mock"
)

type CorporationServiceMock struct {
	mock.Mock
}

func NewCorporationServiceMock() *CorporationServiceMock {
	return &CorporationServiceMock{}
}

func (m *CorporationServiceMock) DoesCorporationExist(corporationID uint) {
	m.Called(corporationID)
}

func (m *CorporationServiceMock) ISCorporationApproved(corporationID uint) bool {
	args := m.Called(corporationID)
	return args.Bool(0)
}

func (m *CorporationServiceMock) GetCorporationCredentials(corporationID uint) corporationdto.CorporationCredentialResponse {
	args := m.Called(corporationID)
	return args.Get(0).(corporationdto.CorporationCredentialResponse)
}

func (m *CorporationServiceMock) CheckApplicantAccess(corporationID, applicantID uint) {
	m.Called(corporationID, applicantID)
}

func (m *CorporationServiceMock) Register(registerInfo corporationdto.RegisterRequest) corporationdto.CorporationCredentialResponse {
	args := m.Called(registerInfo)
	return args.Get(0).(corporationdto.CorporationCredentialResponse)
}

func (m *CorporationServiceMock) UpdateRegister(updateRegisterInfo corporationdto.UpdateRegisterRequest) {
	m.Called(updateRegisterInfo)
}

func (m *CorporationServiceMock) AddCertificateFiles(requestInfo corporationdto.AddCertificatesRequest) {
	m.Called(requestInfo)
}

func (m *CorporationServiceMock) AddContactInfo(contactInfo corporationdto.AddContactInformationRequest) {
	m.Called(contactInfo)
}

func (m *CorporationServiceMock) DeleteContactInfo(contactInfo corporationdto.DeleteContactInformationRequest) {
	m.Called(contactInfo)
}

func (m *CorporationServiceMock) AddAddress(addressInfo corporationdto.AddCorporationAddressRequest) {
	m.Called(addressInfo)
}

func (m *CorporationServiceMock) DeleteAddress(addressInfo corporationdto.DeleteAddressRequest) {
	m.Called(addressInfo)
}

func (m *CorporationServiceMock) GetCorporationDetails(requestInfo corporationdto.CorporationDetailsRequest) corporationdto.CorporationPrivateInfoResponse {
	args := m.Called(requestInfo)
	return args.Get(0).(corporationdto.CorporationPrivateInfoResponse)
}

func (m *CorporationServiceMock) GetContactTypes() []corporationdto.ContactTypeResponse {
	args := m.Called()
	return args.Get(0).([]corporationdto.ContactTypeResponse)
}

func (m *CorporationServiceMock) ChangeLogo(changeLogoRequest corporationdto.ChangeLogoRequest) {
	m.Called(changeLogoRequest)
}

func (m *CorporationServiceMock) GetCorporations(requestInfo corporationdto.CorporationListRequest) []corporationdto.CorporationCredentialResponse {
	args := m.Called(requestInfo)
	return args.Get(0).([]corporationdto.CorporationCredentialResponse)
}
