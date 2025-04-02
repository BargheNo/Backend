package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
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

func (repo *CorporationRepository) DeleteCorporationByCIN(db database.Database, cin string) error {
	return db.GetDB().Where("cin = ?", cin).Delete(&entity.Corporation{}).Error
}

func (repo *CorporationRepository) UpdateCorporation(db database.Database, corporation *entity.Corporation) error {
	return db.GetDB().Save(&corporation).Error
}
