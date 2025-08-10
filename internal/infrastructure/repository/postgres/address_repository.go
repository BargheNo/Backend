package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type AddressRepository struct{}

func NewAddressRepository() *AddressRepository {
	return &AddressRepository{}
}

const (
	queryByOwnerIDAndOwnerType = "owner_id = ? AND owner_type = ?"
)

func (repo *AddressRepository) GetProvinceList(db database.Database) ([]*entity.Province, error) {
	var provinces []*entity.Province
	err := db.GetDB().Find(&provinces).Error
	if err != nil {
		return nil, err
	}
	return provinces, nil
}

func (repo *AddressRepository) GetProvinceCities(db database.Database, provinceID uint) ([]*entity.City, error) {
	var cities []*entity.City
	err := db.GetDB().Where("province_id = ?", provinceID).Find(&cities).Error
	if err != nil {
		return nil, err
	}
	return cities, nil
}

func (repo *AddressRepository) GetProvinceByID(db database.Database, id uint) (*entity.Province, error) {
	var province entity.Province
	result := db.GetDB().First(&province, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &province, nil
}

func (repo *AddressRepository) GetProvinceByName(db database.Database, name string) (*entity.Province, error) {
	var province entity.Province
	result := db.GetDB().Where("name = ?", name).First(&province)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &province, nil
}

func (repo *AddressRepository) GetCityByID(db database.Database, id uint) (*entity.City, error) {
	var city entity.City
	result := db.GetDB().First(&city, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &city, nil
}

func (repo *AddressRepository) GetCityByName(db database.Database, name string) (*entity.City, error) {
	var city entity.City
	result := db.GetDB().Where("name = ?", name).First(&city)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &city, nil
}

func (repo *AddressRepository) GetAddressByID(db database.Database, id uint) (*entity.Address, error) {
	var address entity.Address
	result := db.GetDB().First(&address, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &address, nil
}

func (repo *AddressRepository) GetOwnerAddress(db database.Database, ownerID uint, ownerType string) (*entity.Address, error) {
	var address entity.Address
	result := db.GetDB().Where(queryByOwnerIDAndOwnerType, ownerID, ownerType).First(&address)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &address, nil
}

func (repo *AddressRepository) GetOwnerAddresses(db database.Database, ownerID uint, ownerType string) ([]*entity.Address, error) {
	var addresses []*entity.Address
	err := db.GetDB().Where(queryByOwnerIDAndOwnerType, ownerID, ownerType).Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func (repo *AddressRepository) CreateProvince(db database.Database, province *entity.Province) error {
	return db.GetDB().Create(&province).Error
}

func (repo *AddressRepository) CreateCity(db database.Database, city *entity.City) error {
	return db.GetDB().Create(&city).Error
}

func (repo *AddressRepository) CreateAddress(db database.Database, address *entity.Address) error {
	return db.GetDB().Create(&address).Error
}

func (repo *AddressRepository) DeleteAddress(db database.Database, address *entity.Address) error {
	return db.GetDB().Delete(&address).Error
}
