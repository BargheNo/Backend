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
func (m *CorporationRepositoryMock) FindCorporationByName(db database.Database, name string, status []enum.CorporationStatus) (*entity.Corporation, error) {
	args := m.Called(db, name, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Corporation), args.Error(1)
}

func (m *CorporationRepositoryMock) FindCorporationByRegistrationNumber(db database.Database, registrationNumber string, status []enum.CorporationStatus) (*entity.Corporation, error) {
	args := m.Called(db, registrationNumber, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Corporation), args.Error(1)
}

func (m *CorporationRepositoryMock) FindCorporationByNationalID(db database.Database, nationalID string, status []enum.CorporationStatus) (*entity.Corporation, error) {
	args := m.Called(db, nationalID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Corporation), args.Error(1)
}

func (m *CorporationRepositoryMock) FindCorporationByIBAN(db database.Database, iban string, status []enum.CorporationStatus) (*entity.Corporation, error) {
	args := m.Called(db, iban, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Corporation), args.Error(1)
}

func (m *CorporationRepositoryMock) FindCorporationByCIN(db database.Database, cin string) (*entity.Corporation, error) {
	args := m.Called(db, cin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Corporation), args.Error(1)
}

func (m *CorporationRepositoryMock) FindCorporationByID(db database.Database, id uint) (*entity.Corporation, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Corporation), args.Error(1)
}

func (m *CorporationRepositoryMock) FindCorporationStaff(db database.Database, staffID, corporationID uint) (*entity.CorporationStaff, error) {
	args := m.Called(db, staffID, corporationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CorporationStaff), args.Error(1)
}

func (m *CorporationRepositoryMock) FindContactInformationTypeByID(db database.Database, typeID uint) (*entity.ContactType, error) {
	args := m.Called(db, typeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ContactType), args.Error(1)
}

func (m *CorporationRepositoryMock) FindContactInformationTypeValue(db database.Database, typeID uint, value string) (*entity.ContactInformation, error) {
	args := m.Called(db, typeID, value)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ContactInformation), args.Error(1)
}

func (m *CorporationRepositoryMock) FindContactInformationByID(db database.Database, contactID uint) (*entity.ContactInformation, error) {
	args := m.Called(db, contactID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ContactInformation), args.Error(1)
}

func (m *CorporationRepositoryMock) FindSignatoryByID(db database.Database, signatoryID uint) (*entity.Signatory, error) {
	args := m.Called(db, signatoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Signatory), args.Error(1)
}

func (m *CorporationRepositoryMock) FindCorporationSignatoryByNationalID(db database.Database, corporationID uint, nationalID, position string) (*entity.Signatory, error) {
	args := m.Called(db, corporationID, nationalID, position)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Signatory), args.Error(1)
}

func (m *CorporationRepositoryMock) FindCorporationSignatories(db database.Database, corporationID uint) ([]*entity.Signatory, error) {
	args := m.Called(db, corporationID)
	return args.Get(0).([]*entity.Signatory), args.Error(1)
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

func (m *CorporationRepositoryMock) FindContactTypeByID(db database.Database, contactTypeID uint) (*entity.ContactType, error) {
	args := m.Called(db, contactTypeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ContactType), args.Error(1)
}

func (m *CorporationRepositoryMock) FindContactTypeByName(db database.Database, name string) (*entity.ContactType, error) {
	args := m.Called(db, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ContactType), args.Error(1)
}

func (m *CorporationRepositoryMock) FindContactTypes(db database.Database) ([]*entity.ContactType, error) {
	args := m.Called(db)
	return args.Get(0).([]*entity.ContactType), args.Error(1)
}

func (m *CorporationRepositoryMock) FindUserCorporations(db database.Database, userID uint) ([]*entity.Corporation, error) {
	args := m.Called(db, userID)
	return args.Get(0).([]*entity.Corporation), args.Error(1)
}

func (m *CorporationRepositoryMock) UpdateCorporation(db database.Database, corporation *entity.Corporation) error {
	args := m.Called(db, corporation)
	return args.Error(0)
}

func (m *CorporationRepositoryMock) FindCorporationsByStatus(db database.Database, status []enum.CorporationStatus, opts ...repository.QueryModifier) ([]*entity.Corporation, error) {
	args := m.Called(db, status, opts)
	return args.Get(0).([]*entity.Corporation), args.Error(1)
}

func (m *CorporationRepositoryMock) FindCorporationReviews(db database.Database, corporationID uint, opts ...repository.QueryModifier) ([]*entity.CorporationReview, error) {
	args := m.Called(db, corporationID, opts)
	return args.Get(0).([]*entity.CorporationReview), args.Error(1)
}

func (m *CorporationRepositoryMock) FindContactInformation(db database.Database, corporationID uint) ([]*entity.ContactInformation, error) {
	args := m.Called(db, corporationID)
	return args.Get(0).([]*entity.ContactInformation), args.Error(1)
}

func (m *CorporationRepositoryMock) DeleteCorporationSignatories(db database.Database, corporationID uint) error {
	args := m.Called(db, corporationID)
	return args.Error(0)
}

func (m *CorporationRepositoryMock) DeleteContactInfo(db database.Database, contact *entity.ContactInformation) error {
	args := m.Called(db, contact)
	return args.Error(0)
}

func (m *CorporationRepositoryMock) CreateReview(db database.Database, review *entity.CorporationReview) error {
	args := m.Called(db, review)
	return args.Error(0)
}
