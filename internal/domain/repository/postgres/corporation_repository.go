package repository

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
	FindContactTypeByID(db database.Database, contactTypeID uint) (*entity.ContactType, bool)
	FindContactTypeByName(db database.Database, name string) (*entity.ContactType, bool)
	FindContactTypes(db database.Database) []*entity.ContactType
	FindCorporationByCIN(db database.Database, cin string) (*entity.Corporation, bool)
	FindCorporationByIBAN(db database.Database, iban string, status []enum.CorporationStatus) (*entity.Corporation, bool)
	FindCorporationByID(db database.Database, id uint) (*entity.Corporation, bool)
	FindCorporationByName(db database.Database, name string, status []enum.CorporationStatus) (*entity.Corporation, bool)
	FindCorporationByNationalID(db database.Database, nationalID string, status []enum.CorporationStatus) (*entity.Corporation, bool)
	FindCorporationByRegistrationNumber(db database.Database, registrationNumber string, status []enum.CorporationStatus) (*entity.Corporation, bool)
	FindCorporationStaff(db database.Database, staffID, corporationID uint) (*entity.CorporationStaff, bool)
	FindContactInformationTypeByID(db database.Database, typeID uint) (*entity.ContactType, bool)
	FindContactInformationTypeValue(db database.Database, typeID uint, value string) (*entity.ContactInformation, bool)
	FindContactInformationByID(db database.Database, contactID uint) (*entity.ContactInformation, bool)
	FindSignatoryByID(db database.Database, signatoryID uint) (*entity.Signatory, bool)
	FindCorporationSignatoryByNationalID(db database.Database, corporationID uint, nationalID, position string) (*entity.Signatory, bool)
	FindCorporationSignatories(db database.Database, corporationID uint) []*entity.Signatory
	FindUserCorporations(db database.Database, userID uint) []*entity.Corporation
	UpdateCorporation(db database.Database, corporation *entity.Corporation) error
	FindCorporationsByStatus(db database.Database, status []enum.CorporationStatus, opts ...QueryModifier) []*entity.Corporation
	FindCorporationReviews(db database.Database, corporationID uint, opts ...QueryModifier) []*entity.CorporationReview
	FindContactInformation(db database.Database, corporationID uint) []*entity.ContactInformation
	DeleteCorporationSignatories(db database.Database, corporationID uint) error
	DeleteContactInfo(db database.Database, contact *entity.ContactInformation) error
}
