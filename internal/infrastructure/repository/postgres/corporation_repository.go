package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
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

func (repo *CorporationRepository) FindCorporationByName(db database.Database, name string, status []enum.CorporationStatus) (*entity.Corporation, error) {
	var corporation entity.Corporation
	result := db.GetDB().Where("name = ? AND status IN ?", name, status).First(&corporation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &corporation, nil
}

func (repo *CorporationRepository) FindCorporationByRegistrationNumber(db database.Database, registrationNumber string, status []enum.CorporationStatus) (*entity.Corporation, error) {
	var corporation entity.Corporation
	result := db.GetDB().Where("registration_number = ? AND status IN ?", registrationNumber, status).First(&corporation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &corporation, nil
}

func (repo *CorporationRepository) FindCorporationByNationalID(db database.Database, nationalID string, status []enum.CorporationStatus) (*entity.Corporation, error) {
	var corporation entity.Corporation
	result := db.GetDB().Where("national_id = ? AND status IN ?", nationalID, status).First(&corporation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &corporation, nil
}

func (repo *CorporationRepository) FindCorporationByIBAN(db database.Database, iban string, status []enum.CorporationStatus) (*entity.Corporation, error) {
	var corporation entity.Corporation
	result := db.GetDB().Where("iban = ? AND status IN ?", iban, status).First(&corporation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &corporation, nil
}

func (repo *CorporationRepository) FindCorporationByCIN(db database.Database, cin string) (*entity.Corporation, error) {
	var corporation entity.Corporation
	result := db.GetDB().Where("cin = ?", cin).First(&corporation)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &corporation, nil
}

func (repo *CorporationRepository) FindCorporationByID(db database.Database, id uint) (*entity.Corporation, error) {
	var corporation entity.Corporation
	result := db.GetDB().First(&corporation, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &corporation, nil
}

func (repo *CorporationRepository) FindCorporationStaff(db database.Database, staffID, corporationID uint) (*entity.CorporationStaff, error) {
	var staff entity.CorporationStaff
	result := db.GetDB().Where("user_id = ? AND corporation_ID = ?", staffID, corporationID).First(&staff)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &staff, nil
}

func (repo *CorporationRepository) FindContactInformationTypeByID(db database.Database, typeID uint) (*entity.ContactType, error) {
	var contactType entity.ContactType
	result := db.GetDB().First(&contactType, typeID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &contactType, nil
}

func (repo *CorporationRepository) FindContactInformationTypeValue(db database.Database, typeID uint, value string) (*entity.ContactInformation, error) {
	var contact entity.ContactInformation
	result := db.GetDB().Where("type_id = ? AND value = ?", typeID, value).First(&contact)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &contact, nil
}

func (repo *CorporationRepository) FindContactInformationByID(db database.Database, contactID uint) (*entity.ContactInformation, error) {
	var contact entity.ContactInformation
	result := db.GetDB().First(&contact, contactID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &contact, nil
}

func (repo *CorporationRepository) FindSignatoryByID(db database.Database, signatoryID uint) (*entity.Signatory, error) {
	var signatory entity.Signatory
	result := db.GetDB().First(&signatory, signatoryID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &signatory, nil
}

func (repo *CorporationRepository) FindCorporationSignatoryByNationalID(db database.Database, corporationID uint, nationalID, position string) (*entity.Signatory, error) {
	var signatory entity.Signatory
	result := db.GetDB().Where("corporation_id = ? AND national_card_number = ? AND position = ?", corporationID, nationalID, position).First(&signatory)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &signatory, nil
}

func (repo *CorporationRepository) FindCorporationSignatories(db database.Database, corporationID uint) ([]*entity.Signatory, error) {
	var signatories []*entity.Signatory
	result := db.GetDB().Where("corporation_id = ?", corporationID).Find(&signatories)
	if result.Error != nil {
		return nil, result.Error
	}
	return signatories, nil
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

func (repo *CorporationRepository) FindContactTypeByID(db database.Database, contactTypeID uint) (*entity.ContactType, error) {
	var contactType entity.ContactType
	result := db.GetDB().First(&contactType, contactTypeID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &contactType, nil
}

func (repo *CorporationRepository) FindContactTypeByName(db database.Database, name string) (*entity.ContactType, error) {
	var contactType entity.ContactType
	result := db.GetDB().Where("name = ?", name).First(&contactType)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &contactType, nil
}

func (repo *CorporationRepository) FindContactTypes(db database.Database) ([]*entity.ContactType, error) {
	var types []*entity.ContactType
	err := db.GetDB().Find(&types).Error
	if err != nil {
		return nil, err
	}
	return types, nil
}

func (repo *CorporationRepository) UpdateCorporation(db database.Database, corporation *entity.Corporation) error {
	return db.GetDB().Save(&corporation).Error
}

func (repo *CorporationRepository) FindCorporationsByStatus(db database.Database, status []enum.CorporationStatus, options *postgres.QueryOptions) ([]*entity.Corporation, error) {
	var corporations []*entity.Corporation
	query := db.GetDB().Where("status IN ?", status)
	query = applyQueryOptions(query, options)
	result := query.Find(&corporations)
	if result.Error != nil {
		return nil, result.Error
	}
	return corporations, nil
}

func (repo *CorporationRepository) CountCorporationsByStatus(db database.Database, status []enum.CorporationStatus) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Corporation{}).
		Where("status IN ?", status).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *CorporationRepository) FindCorporationsByQuery(db database.Database, query string, options *postgres.QueryOptions) ([]*entity.Corporation, error) {
	var corporations []*entity.Corporation
	result := db.GetDB().
		Where("name ILIKE ?", "%"+query+"%")

	result = applyQueryOptions(result, options)

	result = result.Find(&corporations)
	if result.Error != nil {
		return nil, result.Error
	}
	return corporations, nil
}

func (repo *CorporationRepository) CountCorporationsByQuery(db database.Database, query string) (int64, error) {
	var count int64
	err := db.GetDB().
		Model(&entity.Corporation{}).
		Where("name ILIKE ?", "%"+query+"%").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *CorporationRepository) FindUserActiveCorporations(db database.Database, userID uint) ([]*entity.Corporation, error) {
	var corporations []*entity.Corporation
	result := db.GetDB().
		Joins("JOIN corporation_staffs ON corporation_staffs.corporation_id = corporations.id").
		Where("corporation_staffs.status = ? AND corporation_staffs.user_id = ?", 1, userID).
		Find(&corporations)

	if result.Error != nil {
		return nil, result.Error
	}
	return corporations, nil
}

func (repo *CorporationRepository) FindContactInformation(db database.Database, corporationID uint) ([]*entity.ContactInformation, error) {
	var contactInfo []*entity.ContactInformation
	result := db.GetDB().Where(queryByCorporationID, corporationID).Find(&contactInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return contactInfo, nil
}

func (repo *CorporationRepository) DeleteCorporationSignatories(db database.Database, corporationID uint) error {
	return db.GetDB().Where(queryByCorporationID, corporationID).Delete(&entity.Signatory{}).Error
}

func (repo *CorporationRepository) DeleteContactInfo(db database.Database, contact *entity.ContactInformation) error {
	return db.GetDB().Delete(contact).Error
}

func (repo *CorporationRepository) FindCorporationReviews(db database.Database, corporationID uint, options *postgres.QueryOptions) ([]*entity.CorporationReview, error) {
	var reviews []*entity.CorporationReview
	query := db.GetDB().Where("corporation_id = ?", corporationID)

	query = applyQueryOptions(query, options)

	result := query.Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}
	return reviews, nil
}

func (repo *CorporationRepository) CreateReview(db database.Database, review *entity.CorporationReview) error {
	return db.GetDB().Create(review).Error
}

func (repo *CorporationRepository) FindStaffRoles(db database.Database, staff *entity.CorporationStaff) error {
	return db.GetDB().Model(staff).Association("Roles").Find(&staff.Roles)
}

func (repo *CorporationRepository) FindStaffByUserIDAndStatus(db database.Database, userID uint, status []enum.StaffStatus) (*entity.CorporationStaff, error) {
	var staff entity.CorporationStaff
	result := db.GetDB().Where("user_id = ? AND status IN ?", userID, status).First(&staff)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &staff, nil
}

func (repo *CorporationRepository) FindCorporationStaffByID(db database.Database, corporationID, staffID uint) (*entity.CorporationStaff, error) {
	var staff entity.CorporationStaff
	result := db.GetDB().Preload("Roles").Where("corporation_id = ?", corporationID).First(&staff, staffID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &staff, nil
}

func (repo *CorporationRepository) FindRolesByIDs(db database.Database, roleIDs []uint, userType enum.UserType) ([]entity.Role, error) {
	var roles []entity.Role
	err := db.GetDB().Where("id IN ? AND user_type = ?", roleIDs, userType).Find(&roles).Error
	return roles, err
}

func (repo *CorporationRepository) CreateStaff(db database.Database, staff *entity.CorporationStaff) error {
	return db.GetDB().Create(staff).Error
}

func (repo *CorporationRepository) UpdateStaff(db database.Database, staff *entity.CorporationStaff) error {
	return db.GetDB().Save(staff).Error
}

func (repo *CorporationRepository) ReplaceStaffRoles(db database.Database, staff *entity.CorporationStaff, roles []entity.Role) error {
	return db.GetDB().Model(&staff).Association("Roles").Replace(roles)
}

func (repo *CorporationRepository) FindRoleByName(db database.Database, name string) (*entity.Role, error) {
	var role entity.Role
	result := db.GetDB().Where("name = ?", name).First(&role)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &role, nil
}

func (repo *CorporationRepository) FindCorporationStaffs(db database.Database, corporationID uint, allowedStatus []enum.StaffStatus, options *postgres.QueryOptions) ([]*entity.CorporationStaff, error) {
	var staffs []*entity.CorporationStaff

	query := db.GetDB().Preload("Roles").Where("corporation_id = ? AND status IN ?", corporationID, allowedStatus)

	query = applyQueryOptions(query, options)

	result := query.Find(&staffs)
	if result.Error != nil {
		return nil, result.Error
	}
	return staffs, nil
}

func (repo *CorporationRepository) CountCorporationStaffs(db database.Database, corporationID uint, allowedStatus []enum.StaffStatus) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.CorporationStaff{}).
		Where("corporation_id = ? AND status IN ?", corporationID, allowedStatus).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *CorporationRepository) FindCorporationStaffByQuery(db database.Database, corporationID uint, allowedStatus []enum.StaffStatus, query string, options *postgres.QueryOptions) ([]*entity.CorporationStaff, error) {
	var staff []*entity.CorporationStaff

	result := db.GetDB().
		Preload("User").
		Preload("Roles").
		Preload("Corporation").
		Joins("LEFT JOIN users ON corporation_staffs.user_id = users.id").
		Where(`
			corporation_staffs.corporation_id = ? AND
			corporation_staffs.status IN ? AND (
			users.first_name ILIKE ? OR
			users.last_name ILIKE ? OR
			users.email ILIKE ? OR
			users.phone ILIKE ?
		)
		`,
			corporationID, allowedStatus, "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")

	result = applyQueryOptions(result, options)
	result = result.Find(&staff)

	if result.Error != nil {
		return nil, result.Error
	}
	return staff, nil
}

func (repo *CorporationRepository) CountCorporationStaffByQuery(db database.Database, corporationID uint, allowedStatus []enum.StaffStatus, query string) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.CorporationStaff{}).
		Joins("LEFT JOIN users ON corporation_staffs.user_id = users.id").
		Where(`
			corporation_staffs.corporation_id = ? AND
			corporation_staffs.status IN ? AND (
			users.first_name ILIKE ? OR
			users.last_name ILIKE ? OR
			users.email ILIKE ? OR
			users.phone ILIKE ?
		)
		`,
			corporationID, allowedStatus, "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}
