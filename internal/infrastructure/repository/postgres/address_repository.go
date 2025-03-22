package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type AddressRepository struct{}

func NewAddressRepository() *AddressRepository {
	return &AddressRepository{}
}

func (repo *AddressRepository) GetProvinceList(db database.Database) []*entity.Province {
	var provinces []*entity.Province
	err := db.GetDB().Find(&provinces).Error
	if err != nil {
		panic(err)
	}
	return provinces
}

func (repo *AddressRepository) GetProvinceCities(db database.Database, provinceID uint) []*entity.City {
	var cities []*entity.City
	err := db.GetDB().Where("province_id = ?", provinceID).Find(&cities).Error
	if err != nil {
		panic(err)
	}
	return cities
}

func (repo *AddressRepository) GetProvinceByID(db database.Database, id uint) (*entity.Province, bool) {
	var province entity.Province
	result := db.GetDB().First(&province, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &province, true
}

func (repo *AddressRepository) GetProvinceByName(db database.Database, name string) (*entity.Province, bool) {
	var province entity.Province
	result := db.GetDB().Where("name = ?", name).First(&province)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &province, true
}

func (repo *AddressRepository) GetCityByID(db database.Database, id uint) (*entity.City, bool) {
	var city entity.City
	result := db.GetDB().First(&city, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &city, true
}

func (repo *AddressRepository) GetCityByName(db database.Database, name string) (*entity.City, bool) {
	var city entity.City
	result := db.GetDB().Where("name = ?", name).First(&city)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &city, true
}

func (repo *AddressRepository) GetAddressByID(db database.Database, id uint) (*entity.Address, bool) {
	var address entity.Address
	result := db.GetDB().First(&address, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &address, true
}

func (repo *AddressRepository) GetOwnerAddresses(db database.Database, ownerID uint, ownerType string) []*entity.Address {
	var addresses []*entity.Address
	err := db.GetDB().Where("owner_id = ? AND owner_type = ?", ownerID, ownerType).Find(&addresses).Error
	if err != nil {
		panic(err)
	}
	return addresses
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
