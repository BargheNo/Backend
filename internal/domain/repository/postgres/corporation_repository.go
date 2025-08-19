package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type CorporationRepository interface {
	CreateCorporation(db database.Database, corporation *entity.Corporation) error
	CreateCorporationStaff(db database.Database, staff *entity.CorporationStaff) error
	CreateSignatory(db database.Database, signatory *entity.Signatory) error
	CreateContactInformation(db database.Database, contact *entity.ContactInformation) error
	CreateContactType(db database.Database, contactType *entity.ContactType) error
	FindContactTypeByID(db database.Database, contactTypeID uint) (*entity.ContactType, error)
	FindContactTypeByName(db database.Database, name string) (*entity.ContactType, error)
	FindContactTypes(db database.Database) ([]*entity.ContactType, error)
	FindCorporationByCIN(db database.Database, cin string) (*entity.Corporation, error)
	FindCorporationByIBAN(db database.Database, iban string, status []enum.CorporationStatus) (*entity.Corporation, error)
	FindCorporationByID(db database.Database, id uint) (*entity.Corporation, error)
	FindCorporationByName(db database.Database, name string, status []enum.CorporationStatus) (*entity.Corporation, error)
	FindCorporationByNationalID(db database.Database, nationalID string, status []enum.CorporationStatus) (*entity.Corporation, error)
	FindCorporationByRegistrationNumber(db database.Database, registrationNumber string, status []enum.CorporationStatus) (*entity.Corporation, error)
	FindCorporationStaff(db database.Database, staffID, corporationID uint) (*entity.CorporationStaff, error)
	FindContactInformationTypeByID(db database.Database, typeID uint) (*entity.ContactType, error)
	FindContactInformationTypeValue(db database.Database, typeID uint, value string) (*entity.ContactInformation, error)
	FindContactInformationByID(db database.Database, contactID uint) (*entity.ContactInformation, error)
	FindSignatoryByID(db database.Database, signatoryID uint) (*entity.Signatory, error)
	FindCorporationSignatoryByNationalID(db database.Database, corporationID uint, nationalID, position string) (*entity.Signatory, error)
	FindCorporationSignatories(db database.Database, corporationID uint) ([]*entity.Signatory, error)
	FindUserActiveCorporations(db database.Database, userID uint) ([]*entity.Corporation, error)
	UpdateCorporation(db database.Database, corporation *entity.Corporation) error
	FindCorporationsByStatus(db database.Database, status []enum.CorporationStatus, options *QueryOptions) ([]*entity.Corporation, error)
	FindCorporationsByQuery(db database.Database, query string, options *QueryOptions) ([]*entity.Corporation, error)
	CountCorporationsByQuery(db database.Database, query string) (int64, error)
	CountCorporationsByStatus(db database.Database, status []enum.CorporationStatus) (int64, error)
	FindCorporationReviews(db database.Database, corporationID uint, options *QueryOptions) ([]*entity.CorporationReview, error)
	FindContactInformation(db database.Database, corporationID uint) ([]*entity.ContactInformation, error)
	DeleteCorporationSignatories(db database.Database, corporationID uint) error
	DeleteContactInfo(db database.Database, contact *entity.ContactInformation) error
	CreateReview(db database.Database, review *entity.CorporationReview) error
	FindStaffRoles(db database.Database, staff *entity.CorporationStaff) error
	FindStaffByUserIDAndStatus(db database.Database, userID uint, status []enum.StaffStatus) (*entity.CorporationStaff, error)
	FindRolesByIDs(db database.Database, roleIDs []uint, userType enum.UserType) ([]entity.Role, error)
	CreateStaff(db database.Database, staff *entity.CorporationStaff) error
	ReplaceStaffRoles(db database.Database, staff *entity.CorporationStaff, roles []entity.Role) error
	FindRoleByName(db database.Database, name string) (*entity.Role, error)
	FindCorporationStaffByID(db database.Database, corporationID, staffID uint) (*entity.CorporationStaff, error)
	UpdateStaff(db database.Database, staff *entity.CorporationStaff) error
	FindCorporationStaffs(db database.Database, corporationID uint, allowedStatus []enum.StaffStatus, options *QueryOptions) ([]*entity.CorporationStaff, error)
	CountCorporationStaffs(db database.Database, corporationID uint, allowedStatus []enum.StaffStatus) (int64, error)
	FindCorporationStaffByQuery(db database.Database, corporationID uint, allowedStatus []enum.StaffStatus, query string, options *QueryOptions) ([]*entity.CorporationStaff, error)
	CountCorporationStaffByQuery(db database.Database, corporationID uint, allowedStatus []enum.StaffStatus, query string) (int64, error)
}
