package mocks

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type CorporationRepositoryMock struct {
	mock.Mock
}

func NewCorporationRepositoryMock() *CorporationRepositoryMock {
	return &CorporationRepositoryMock{}
}

func (c *CorporationRepositoryMock) CreateCorporation(db database.Database, corporation *entity.Corporation) error {
	args := c.Called(db, corporation)
	return args.Error(0)
}

func (c *CorporationRepositoryMock) CreateCorporationStaff(db database.Database, staff *entity.CorporationStaff) error {
	args := c.Called(db, staff)
	return args.Error(0)
}

func (c *CorporationRepositoryMock) CreateSignatory(db database.Database, signatory *entity.Signatory) error {
	args := c.Called(db, signatory)
	return args.Error(0)
}

func (c *CorporationRepositoryMock) CreateContactInformation(db database.Database, contact *entity.ContactInformation) error {
	args := c.Called(db, contact)
	return args.Error(0)
}

func (c *CorporationRepositoryMock) CreateContactType(db database.Database, contactType *entity.ContactType) error {
	args := c.Called(db, contactType)
	return args.Error(0)
}

func (c *CorporationRepositoryMock) FindContactTypeByID(db database.Database, contactTypeID uint) (*entity.ContactType, error) {
	args := c.Called(db, contactTypeID)
	return args.Get(0).(*entity.ContactType), args.Error(1)
}

func (c *CorporationRepositoryMock) FindContactTypeByName(db database.Database, name string) (*entity.ContactType, error) {
	args := c.Called(db, name)
	return args.Get(0).(*entity.ContactType), args.Error(1)
}

func (c *CorporationRepositoryMock) FindContactTypes(db database.Database) ([]*entity.ContactType, error) {
	args := c.Called(db)
	return args.Get(0).([]*entity.ContactType), args.Error(1)
}

func (c *CorporationRepositoryMock) FindCorporationByCIN(db database.Database, cin string) (*entity.Corporation, error) {
	args := c.Called(db, cin)
	return args.Get(0).(*entity.Corporation), args.Error(1)
}

func (c *CorporationRepositoryMock) FindCorporationByIBAN(db database.Database, iban string, status []enum.CorporationStatus) (*entity.Corporation, error) {
	args := c.Called(db, iban, status)
	return args.Get(0).(*entity.Corporation), args.Error(1)
}

func (c *CorporationRepositoryMock) FindCorporationByID(db database.Database, id uint) (*entity.Corporation, error) {
	args := c.Called(db, id)
	return args.Get(0).(*entity.Corporation), args.Error(1)
}

func (c *CorporationRepositoryMock) FindCorporationByName(db database.Database, name string, status []enum.CorporationStatus) (*entity.Corporation, error) {
	args := c.Called(db, name, status)
	return args.Get(0).(*entity.Corporation), args.Error(1)
}

func (c *CorporationRepositoryMock) FindCorporationByNationalID(db database.Database, nationalID string, status []enum.CorporationStatus) (*entity.Corporation, error) {
	args := c.Called(db, nationalID, status)
	return args.Get(0).(*entity.Corporation), args.Error(1)
}

func (c *CorporationRepositoryMock) FindCorporationByRegistrationNumber(db database.Database, registrationNumber string, status []enum.CorporationStatus) (*entity.Corporation, error) {
	args := c.Called(db, registrationNumber, status)
	return args.Get(0).(*entity.Corporation), args.Error(1)
}

func (c *CorporationRepositoryMock) FindCorporationStaff(db database.Database, staffID, corporationID uint) (*entity.CorporationStaff, error) {
	args := c.Called(db, staffID, corporationID)
	return args.Get(0).(*entity.CorporationStaff), args.Error(1)
}

func (c *CorporationRepositoryMock) FindContactInformationTypeByID(db database.Database, typeID uint) (*entity.ContactType, error) {
	args := c.Called(db, typeID)
	return args.Get(0).(*entity.ContactType), args.Error(1)
}

func (c *CorporationRepositoryMock) FindContactInformationTypeValue(db database.Database, typeID uint, value string) (*entity.ContactInformation, error) {
	args := c.Called(db, typeID, value)
	return args.Get(0).(*entity.ContactInformation), args.Error(1)
}

func (c *CorporationRepositoryMock) FindContactInformationByID(db database.Database, contactID uint) (*entity.ContactInformation, error) {
	args := c.Called(db, contactID)
	return args.Get(0).(*entity.ContactInformation), args.Error(1)
}

func (c *CorporationRepositoryMock) FindSignatoryByID(db database.Database, signatoryID uint) (*entity.Signatory, error) {
	args := c.Called(db, signatoryID)
	return args.Get(0).(*entity.Signatory), args.Error(1)
}

func (c *CorporationRepositoryMock) FindCorporationSignatoryByNationalID(db database.Database, corporationID uint, nationalID, position string) (*entity.Signatory, error) {
	args := c.Called(db, corporationID, nationalID, position)
	return args.Get(0).(*entity.Signatory), args.Error(1)
}

func (c *CorporationRepositoryMock) FindCorporationSignatories(db database.Database, corporationID uint) ([]*entity.Signatory, error) {
	args := c.Called(db, corporationID)
	return args.Get(0).([]*entity.Signatory), args.Error(1)
}

func (c *CorporationRepositoryMock) FindUserActiveCorporations(db database.Database, userID uint) ([]*entity.Corporation, error) {
	args := c.Called(db, userID)
	return args.Get(0).([]*entity.Corporation), args.Error(1)
}

func (c *CorporationRepositoryMock) UpdateCorporation(db database.Database, corporation *entity.Corporation) error {
	args := c.Called(db, corporation)
	return args.Error(0)
}

func (c *CorporationRepositoryMock) FindCorporationsByStatus(db database.Database, status []enum.CorporationStatus, options *postgres.QueryOptions) ([]*entity.Corporation, error) {
	args := c.Called(db, status, options)
	return args.Get(0).([]*entity.Corporation), args.Error(1)
}

func (c *CorporationRepositoryMock) FindCorporationsByQuery(db database.Database, query string, options *postgres.QueryOptions) ([]*entity.Corporation, error) {
	args := c.Called(db, query, options)
	return args.Get(0).([]*entity.Corporation), args.Error(1)
}

func (c *CorporationRepositoryMock) CountCorporationsByQuery(db database.Database, query string) (int64, error) {
	args := c.Called(db, query)
	return args.Get(0).(int64), args.Error(1)
}

func (c *CorporationRepositoryMock) CountCorporationsByStatus(db database.Database, status []enum.CorporationStatus) (int64, error) {
	args := c.Called(db, status)
	return args.Get(0).(int64), args.Error(1)
}

func (c *CorporationRepositoryMock) FindCorporationReviews(db database.Database, corporationID uint, options *postgres.QueryOptions) ([]*entity.CorporationReview, error) {
	args := c.Called(db, corporationID, options)
	return args.Get(0).([]*entity.CorporationReview), args.Error(1)
}

func (c *CorporationRepositoryMock) FindContactInformation(db database.Database, corporationID uint) ([]*entity.ContactInformation, error) {
	args := c.Called(db, corporationID)
	return args.Get(0).([]*entity.ContactInformation), args.Error(1)
}

func (c *CorporationRepositoryMock) DeleteCorporationSignatories(db database.Database, corporationID uint) error {
	args := c.Called(db, corporationID)
	return args.Error(0)
}

func (c *CorporationRepositoryMock) DeleteContactInfo(db database.Database, contact *entity.ContactInformation) error {
	args := c.Called(db, contact)
	return args.Error(0)
}

func (c *CorporationRepositoryMock) CreateReview(db database.Database, review *entity.CorporationReview) error {
	args := c.Called(db, review)
	return args.Error(0)
}

func (c *CorporationRepositoryMock) FindStaffRoles(db database.Database, staff *entity.CorporationStaff) error {
	args := c.Called(db, staff)
	return args.Error(0)
}

func (c *CorporationRepositoryMock) FindStaffByUserIDAndStatus(db database.Database, userID uint, status []enum.StaffStatus) (*entity.CorporationStaff, error) {
	args := c.Called(db, userID, status)
	return args.Get(0).(*entity.CorporationStaff), args.Error(1)
}

func (c *CorporationRepositoryMock) FindRolesByIDs(db database.Database, roleIDs []uint, userType enum.UserType) ([]entity.Role, error) {
	args := c.Called(db, roleIDs, userType)
	return args.Get(0).([]entity.Role), args.Error(1)
}

func (c *CorporationRepositoryMock) CreateStaff(db database.Database, staff *entity.CorporationStaff) error {
	args := c.Called(db, staff)
	return args.Error(0)
}

func (c *CorporationRepositoryMock) ReplaceStaffRoles(db database.Database, staff *entity.CorporationStaff, roles []entity.Role) error {
	args := c.Called(db, staff, roles)
	return args.Error(0)
}

func (c *CorporationRepositoryMock) FindRoleByName(db database.Database, name string) (*entity.Role, error) {
	args := c.Called(db, name)
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (c *CorporationRepositoryMock) FindCorporationStaffByID(db database.Database, corporationID, staffID uint) (*entity.CorporationStaff, error) {
	args := c.Called(db, corporationID, staffID)
	return args.Get(0).(*entity.CorporationStaff), args.Error(1)
}

func (c *CorporationRepositoryMock) UpdateStaff(db database.Database, staff *entity.CorporationStaff) error {
	args := c.Called(db, staff)
	return args.Error(0)
}

func (c *CorporationRepositoryMock) FindCorporationStaffs(db database.Database, corporationID uint, allowedStatus []enum.StaffStatus, options *postgres.QueryOptions) ([]*entity.CorporationStaff, error) {
	args := c.Called(db, corporationID, allowedStatus, options)
	return args.Get(0).([]*entity.CorporationStaff), args.Error(1)
}

func (c *CorporationRepositoryMock) CountCorporationStaffs(db database.Database, corporationID uint, allowedStatus []enum.StaffStatus) (int64, error) {
	args := c.Called(db, corporationID, allowedStatus)
	return args.Get(0).(int64), args.Error(1)
}

func (c *CorporationRepositoryMock) FindCorporationStaffByQuery(db database.Database, corporationID uint, allowedStatus []enum.StaffStatus, query string, options *postgres.QueryOptions) ([]*entity.CorporationStaff, error) {
	args := c.Called(db, corporationID, allowedStatus, query, options)
	return args.Get(0).([]*entity.CorporationStaff), args.Error(1)
}

func (c *CorporationRepositoryMock) CountCorporationStaffByQuery(db database.Database, corporationID uint, allowedStatus []enum.StaffStatus, query string) (int64, error) {
	args := c.Called(db, corporationID, allowedStatus, query)
	return args.Get(0).(int64), args.Error(1)
}
