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
	DeleteCorporationByCIN(db database.Database, cin string) error
	FindCorporationByCIN(db database.Database, cin string) (*entity.Corporation, bool)
	FindCorporationByIBAN(db database.Database, iban string, status []enum.CorporationStatus) (*entity.Corporation, bool)
	FindCorporationByID(db database.Database, id uint) (*entity.Corporation, bool)
	FindCorporationByName(db database.Database, name string, status []enum.CorporationStatus) (*entity.Corporation, bool)
	FindCorporationByNationalID(db database.Database, nationalID string, status []enum.CorporationStatus) (*entity.Corporation, bool)
	FindCorporationByRegistrationNumber(db database.Database, registrationNumber string, status []enum.CorporationStatus) (*entity.Corporation, bool)
	FindCorporationStaff(db database.Database, staffID, corporationID uint) (*entity.CorporationStaff, bool)
	FindContactInformationTypeByID(db database.Database, typeID uint) (*entity.ContactType, bool)
	FindContactInformationTypeValue(db database.Database, typeID uint, value string) (*entity.ContactInformation, bool)
	UpdateCorporation(db database.Database, corporation *entity.Corporation) error
}
