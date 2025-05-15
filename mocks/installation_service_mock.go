package mocks

import (
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type InstallationServiceMock struct {
	mock.Mock
}

func NewInstallationServiceMock() *InstallationServiceMock {
	return &InstallationServiceMock{}
}

func (m *InstallationServiceMock) GetInstallationRequestModel(requestID uint) *entity.InstallationRequest {
	args := m.Called(requestID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*entity.InstallationRequest)
}

func (m *InstallationServiceMock) CreateInstallationRequest(requestInfo installationdto.NewInstallationRequest) {
	m.Called(requestInfo)
}

func (m *InstallationServiceMock) GetOwnerInstallationRequests(listInfo installationdto.InstallationListRequest) []installationdto.OwnerRequestsResponse {
	args := m.Called(listInfo)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]installationdto.OwnerRequestsResponse)
}

func (m *InstallationServiceMock) GetInstallationRequest(requestID uint) installationdto.RequestDetailsResponse {
	args := m.Called(requestID)
	return args.Get(0).(installationdto.RequestDetailsResponse)
}

func (m *InstallationServiceMock) GetOwnerInstallationRequest(requestInfo installationdto.GetOwnerRequest) installationdto.OwnerRequestsResponse {
	args := m.Called(requestInfo)
	return args.Get(0).(installationdto.OwnerRequestsResponse)
}

func (m *InstallationServiceMock) GetInstallationRequests(listInfo installationdto.InstallationListRequest) []installationdto.RequestDetailsResponse {
	args := m.Called(listInfo)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]installationdto.RequestDetailsResponse)
}

func (m *InstallationServiceMock) AddPanel(panelInfo installationdto.AddPanelRequest) {
	m.Called(panelInfo)
}

func (m *InstallationServiceMock) GetCorporationPanels(listInfo installationdto.CorporationPanelListRequest) []installationdto.CorporationPanelResponse {
	args := m.Called(listInfo)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]installationdto.CorporationPanelResponse)
}

func (m *InstallationServiceMock) GetCustomerPanels(listInfo installationdto.CustomerPanelListRequest) []installationdto.CustomerPanelResponse {
	args := m.Called(listInfo)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]installationdto.CustomerPanelResponse)
}

func (m *InstallationServiceMock) GetPanelByID(panelID uint) installationdto.PanelResponse {
	args := m.Called(panelID)
	return args.Get(0).(installationdto.PanelResponse)
}

func (m *InstallationServiceMock) GetCustomerPanelByID(panelID uint) installationdto.CustomerPanelResponse {
	args := m.Called(panelID)
	return args.Get(0).(installationdto.CustomerPanelResponse)
}

func (m *InstallationServiceMock) GetCorporationPanelByID(panelID uint) installationdto.CorporationPanelResponse {
	args := m.Called(panelID)
	return args.Get(0).(installationdto.CorporationPanelResponse)
}
