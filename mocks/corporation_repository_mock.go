package mocks

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type CorporationRepositoryMock struct {
	mock.Mock
}

func NewCorporationRepositoryMock() *CorporationRepositoryMock {
	return &CorporationRepositoryMock{}
}

// CorporationRepository interface methods
func (m *CorporationRepositoryMock) FindCorporationByName(db database.Database, name string, status []enum.CorporationStatus) (*entity.Corporation, bool) {
	args := m.Called(db, name, status)
	return args.Get(0).(*entity.Corporation), args.Bool(1)
}

func (m *CorporationRepositoryMock) FindCorporationByRegistrationNumber(db database.Database, registrationNumber string, status []enum.CorporationStatus) (*entity.Corporation, bool) {
	args := m.Called(db, registrationNumber, status)
	return args.Get(0).(*entity.Corporation), args.Bool(1)
}

func (m *CorporationRepositoryMock) FindCorporationByNationalID(db database.Database, nationalID string, status []enum.CorporationStatus) (*entity.Corporation, bool) {
	args := m.Called(db, nationalID, status)
	return args.Get(0).(*entity.Corporation), args.Bool(1)
}

func (m *CorporationRepositoryMock) FindCorporationByIBAN(db database.Database, iban string, status []enum.CorporationStatus) (*entity.Corporation, bool) {
	args := m.Called(db, iban, status)
	return args.Get(0).(*entity.Corporation), args.Bool(1)
}

func (m *CorporationRepositoryMock) FindCorporationByCIN(db database.Database, cin string) (*entity.Corporation, bool) {
	args := m.Called(db, cin)
	return args.Get(0).(*entity.Corporation), args.Bool(1)
}

func (m *CorporationRepositoryMock) FindCorporationByID(db database.Database, id uint) (*entity.Corporation, bool) {
	args := m.Called(db, id)
	return args.Get(0).(*entity.Corporation), args.Bool(1)
}

func (m *CorporationRepositoryMock) FindCorporationStaff(db database.Database, staffID, corporationID uint) (*entity.CorporationStaff, bool) {
	args := m.Called(db, staffID, corporationID)
	return args.Get(0).(*entity.CorporationStaff), args.Bool(1)
}

func (m *CorporationRepositoryMock) FindContactInformationTypeByID(db database.Database, typeID uint) (*entity.ContactType, bool) {
	args := m.Called(db, typeID)
	return args.Get(0).(*entity.ContactType), args.Bool(1)
}

func (m *CorporationRepositoryMock) FindContactInformationTypeValue(db database.Database, typeID uint, value string) (*entity.ContactInformation, bool) {
	args := m.Called(db, typeID, value)
	return args.Get(0).(*entity.ContactInformation), args.Bool(1)
}

func (m *CorporationRepositoryMock) FindContactInformationByID(db database.Database, contactID uint) (*entity.ContactInformation, bool) {
	args := m.Called(db, contactID)
	return args.Get(0).(*entity.ContactInformation), args.Bool(1)
}

func (m *CorporationRepositoryMock) FindSignatoryByID(db database.Database, signatoryID uint) (*entity.Signatory, bool) {
	args := m.Called(db, signatoryID)
	return args.Get(0).(*entity.Signatory), args.Bool(1)
}

func (m *CorporationRepositoryMock) FindCorporationSignatoryByNationalID(db database.Database, corporationID uint, nationalID, position string) (*entity.Signatory, bool) {
	args := m.Called(db, corporationID, nationalID, position)
	return args.Get(0).(*entity.Signatory), args.Bool(1)
}

func (m *CorporationRepositoryMock) FindCorporationSignatories(db database.Database, corporationID uint) []*entity.Signatory {
	args := m.Called(db, corporationID)
	return args.Get(0).([]*entity.Signatory)
}

func (m *CorporationRepositoryMock) CreateCorporation(db database.Database, corporation *entity.Corporation) error {
	args := m.Called(db, corporation)
	return args.Error(0)
}

func (m *CorporationRepositoryMock) CreateCorporationStaff(db database.Database, staff *entity.CorporationStaff) error {
	args := m.Called(db, staff)
	return args.Error(0)
}

func (m *CorporationRepositoryMock) CreateSignatory(db database.Database, signatory *entity.Signatory) error {
	args := m.Called(db, signatory)
	return args.Error(0)
}

func (m *CorporationRepositoryMock) CreateContactInformation(db database.Database, contact *entity.ContactInformation) error {
	args := m.Called(db, contact)
	return args.Error(0)
}

func (m *CorporationRepositoryMock) CreateContactType(db database.Database, contactType *entity.ContactType) error {
	args := m.Called(db, contactType)
	return args.Error(0)
}

func (m *CorporationRepositoryMock) FindContactTypeByID(db database.Database, contactTypeID uint) (*entity.ContactType, bool) {
	args := m.Called(db, contactTypeID)
	return args.Get(0).(*entity.ContactType), args.Bool(1)
}

func (m *CorporationRepositoryMock) FindContactTypeByName(db database.Database, name string) (*entity.ContactType, bool) {
	args := m.Called(db, name)
	return args.Get(0).(*entity.ContactType), args.Bool(1)
}

func (m *CorporationRepositoryMock) FindContactTypes(db database.Database) []*entity.ContactType {
	args := m.Called(db)
	return args.Get(0).([]*entity.ContactType)
}

func (m *CorporationRepositoryMock) UpdateCorporation(db database.Database, corporation *entity.Corporation) error {
	args := m.Called(db, corporation)
	return args.Error(0)
}

func (m *CorporationRepositoryMock) FindCorporationsByStatus(db database.Database, status []enum.CorporationStatus, opts ...repository.QueryModifier) []*entity.Corporation {
	args := m.Called(db, status, opts)
	return args.Get(0).([]*entity.Corporation)
}

func (m *CorporationRepositoryMock) FindContactInformation(db database.Database, corporationID uint) []*entity.ContactInformation {
	args := m.Called(db, corporationID)
	return args.Get(0).([]*entity.ContactInformation)
}

func (m *CorporationRepositoryMock) DeleteCorporationSignatories(db database.Database, corporationID uint) error {
	args := m.Called(db, corporationID)
	return args.Error(0)
}

func (m *CorporationRepositoryMock) DeleteContactInfo(db database.Database, contact *entity.ContactInformation) error {
	args := m.Called(db, contact)
	return args.Error(0)
}
