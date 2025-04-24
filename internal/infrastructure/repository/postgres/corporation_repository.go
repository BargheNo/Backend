package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

const (
	queryByCorporationID string = "corporation_id"
)

type CorporationRepository struct{}

func NewCorporationRepository() *CorporationRepository {
	return &CorporationRepository{}
}

func (repo *CorporationRepository) FindCorporationByName(db database.Database, name string, status []enum.CorporationStatus) (*entity.Corporation, bool) {
	var corporation entity.Corporation
	result := db.GetDB().Where("name = ? AND status IN ?", name, status).First(&corporation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &corporation, true
}

func (repo *CorporationRepository) FindCorporationByRegistrationNumber(db database.Database, registrationNumber string, status []enum.CorporationStatus) (*entity.Corporation, bool) {
	var corporation entity.Corporation
	result := db.GetDB().Where("registration_number = ? AND status IN ?", registrationNumber, status).First(&corporation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &corporation, true
}

func (repo *CorporationRepository) FindCorporationByNationalID(db database.Database, nationalID string, status []enum.CorporationStatus) (*entity.Corporation, bool) {
	var corporation entity.Corporation
	result := db.GetDB().Where("national_id = ? AND status IN ?", nationalID, status).First(&corporation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &corporation, true
}

func (repo *CorporationRepository) FindCorporationByIBAN(db database.Database, iban string, status []enum.CorporationStatus) (*entity.Corporation, bool) {
	var corporation entity.Corporation
	result := db.GetDB().Where("iban = ? AND status IN ?", iban, status).First(&corporation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &corporation, true
}

func (repo *CorporationRepository) FindCorporationByCIN(db database.Database, cin string) (*entity.Corporation, bool) {
	var corporation entity.Corporation
	result := db.GetDB().Where("cin = ?", cin).First(&corporation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &corporation, true
}

func (repo *CorporationRepository) FindCorporationByID(db database.Database, id uint) (*entity.Corporation, bool) {
	var corporation entity.Corporation
	result := db.GetDB().First(&corporation, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &corporation, true
}

func (repo *CorporationRepository) FindCorporationStaff(db database.Database, staffID, corporationID uint) (*entity.CorporationStaff, bool) {
	var staff entity.CorporationStaff
	result := db.GetDB().Where("staff_id = ? AND corporation_ID = ?", staffID, corporationID).First(&staff)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &staff, true
}

func (repo *CorporationRepository) FindContactInformationTypeByID(db database.Database, typeID uint) (*entity.ContactType, bool) {
	var contactType entity.ContactType
	result := db.GetDB().First(&contactType, typeID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &contactType, true
}

func (repo *CorporationRepository) FindContactInformationTypeValue(db database.Database, typeID uint, value string) (*entity.ContactInformation, bool) {
	var contact entity.ContactInformation
	result := db.GetDB().Where("type_id = ? AND value = ?", typeID, value).First(&contact)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &contact, true
}

func (repo *CorporationRepository) FindContactInformationByID(db database.Database, contactID uint) (*entity.ContactInformation, bool) {
	var contact entity.ContactInformation
	result := db.GetDB().First(&contact, contactID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &contact, true
}

func (repo *CorporationRepository) FindSignatoryByID(db database.Database, signatoryID uint) (*entity.Signatory, bool) {
	var signatory entity.Signatory
	result := db.GetDB().First(&signatory, signatoryID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &signatory, true
}

func (repo *CorporationRepository) FindCorporationSignatoryByNationalID(db database.Database, corporationID uint, nationalID, position string) (*entity.Signatory, bool) {
	var signatory entity.Signatory
	result := db.GetDB().Where("corporation_id = ? AND national_card_number = ? AND position = ?", corporationID, nationalID, position).First(&signatory)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &signatory, true
}

func (repo *CorporationRepository) CreateCorporation(db database.Database, corporation *entity.Corporation) error {
	return db.GetDB().Create(&corporation).Error
}

func (repo *CorporationRepository) CreateCorporationStaff(db database.Database, staff *entity.CorporationStaff) error {
	return db.GetDB().Create(&staff).Error
}

func (repo *CorporationRepository) CreateSignatory(db database.Database, signatory *entity.Signatory) error {
	return db.GetDB().Create(&signatory).Error
}

func (repo *CorporationRepository) CreateContactInformation(db database.Database, contact *entity.ContactInformation) error {
	return db.GetDB().Create(&contact).Error
}

func (repo *CorporationRepository) CreateContactType(db database.Database, contactType *entity.ContactType) error {
	return db.GetDB().Create(&contactType).Error
}

func (repo *CorporationRepository) FindContactTypeByID(db database.Database, contactTypeID uint) (*entity.ContactType, bool) {
	var contactType entity.ContactType
	result := db.GetDB().First(&contactType, contactTypeID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &contactType, true
}

func (repo *CorporationRepository) FindContactTypeByName(db database.Database, name string) (*entity.ContactType, bool) {
	var contactType entity.ContactType
	result := db.GetDB().Where("name = ?", name).First(&contactType)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &contactType, true
}

func (repo *CorporationRepository) FindContactTypes(db database.Database) []*entity.ContactType {
	var types []*entity.ContactType
	err := db.GetDB().Find(&types).Error
	if err != nil {
		panic(err)
	}
	return types
}

func (repo *CorporationRepository) UpdateCorporation(db database.Database, corporation *entity.Corporation) error {
	return db.GetDB().Save(&corporation).Error
}

func (repo *CorporationRepository) FindCorporationByStatus(db database.Database, status []enum.CorporationStatus, opts ...repository.QueryModifier) []*entity.Corporation {
	var corporations []*entity.Corporation
	query := db.GetDB().Where("status IN ?", status)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}
	result := query.Find(&corporations)
	if result.Error != nil {
		panic(result.Error)
	}
	return corporations
}

func (repo *CorporationRepository) FindContactInformation(db database.Database, corporationID uint) []*entity.ContactInformation {
	var contactInfo []*entity.ContactInformation
	result := db.GetDB().Where(queryByCorporationID, corporationID).Find(&contactInfo)
	if result.Error != nil {
		panic(result.Error)
	}
	return contactInfo
}

func (repo *CorporationRepository) DeleteCorporationSignatories(db database.Database, corporationID uint) error {
	return db.GetDB().Where(queryByCorporationID, corporationID).Delete(&entity.Signatory{}).Error
}

func (repo *CorporationRepository) DeleteContactInfo(db database.Database, contact *entity.ContactInformation) error {
	return db.GetDB().Delete(contact).Error
}
