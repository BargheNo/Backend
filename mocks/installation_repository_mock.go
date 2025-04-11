package mocks

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type InstallationRepositoryMock struct {
	mock.Mock
}

func NewInstallationRepositoryMock() *InstallationRepositoryMock {
	return &InstallationRepositoryMock{}
}

func (repo *InstallationRepositoryMock) FindRequestByStatus(db database.Database, status []enum.InstallationRequestStatus, modifiers ...repository.QueryModifier) []*entity.InstallationRequest {
	args := repo.Called(status)
	return args.Get(0).([]*entity.InstallationRequest)
}

func (repo *InstallationRepositoryMock) FindRequestByID(db database.Database, requestID uint) (*entity.InstallationRequest, bool) {
	args := repo.Called(requestID)
	return args.Get(0).(*entity.InstallationRequest), args.Bool(1)
}

func (repo *InstallationRepositoryMock) FindOwnerRequests(db database.Database, ownerID uint, status []enum.InstallationRequestStatus, modifiers ...repository.QueryModifier) []*entity.InstallationRequest {
	args := repo.Called(ownerID, status)
	return args.Get(0).([]*entity.InstallationRequest)
}

func (repo *InstallationRepositoryMock) FindOwnerRequestByName(db database.Database, ownerID uint, status []enum.InstallationRequestStatus, name string) (*entity.InstallationRequest, bool) {
	args := repo.Called(ownerID, status, name)
	return args.Get(0).(*entity.InstallationRequest), args.Bool(1)
}

func (repo *InstallationRepositoryMock) CreateRequest(db database.Database, request *entity.InstallationRequest) error {
	args := repo.Called(request)
	return args.Error(0)
}

func (repo *InstallationRepositoryMock) CreatePanel(db database.Database, panel *entity.Panel) error {
	args := repo.Called(panel)
	return args.Error(0)
}

func (repo *InstallationRepositoryMock) FindCorporationPanels(db database.Database, corporationID uint, modifiers ...repository.QueryModifier) []*entity.Panel {
	args := repo.Called(corporationID)
	return args.Get(0).([]*entity.Panel)
}

func (repo *InstallationRepositoryMock) FindCustomerPanels(db database.Database, customerID uint, modifiers ...repository.QueryModifier) []*entity.Panel {
	args := repo.Called(customerID)
	return args.Get(0).([]*entity.Panel)
}

func (repo *InstallationRepositoryMock) FindPanelByNameAndCustomerID(db database.Database, panelName string, customerID uint) (*entity.Panel, bool) {
	args := repo.Called(panelName, customerID)
	return args.Get(0).(*entity.Panel), args.Bool(1)
}

func (repo *InstallationRepositoryMock) FindPanelByID(db database.Database, panelID uint) (*entity.Panel, bool) {
	args := repo.Called(panelID)
	return args.Get(0).(*entity.Panel), args.Bool(1)
}
