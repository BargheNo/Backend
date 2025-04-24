package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type AddressRepository interface {
	CreateAddress(db database.Database, address *entity.Address) error
	CreateCity(db database.Database, city *entity.City) error
	CreateProvince(db database.Database, province *entity.Province) error
	DeleteAddress(db database.Database, address *entity.Address) error
	GetAddressByID(db database.Database, id uint) (*entity.Address, bool)
	GetCityByID(db database.Database, id uint) (*entity.City, bool)
	GetCityByName(db database.Database, name string) (*entity.City, bool)
	GetOwnerAddress(db database.Database, ownerID uint, ownerType string) (*entity.Address, bool)
	GetOwnerAddresses(db database.Database, ownerID uint, ownerType string) []*entity.Address
	GetProvinceByID(db database.Database, id uint) (*entity.Province, bool)
	GetProvinceByName(db database.Database, name string) (*entity.Province, bool)
	GetProvinceCities(db database.Database, provinceID uint) []*entity.City
	GetProvinceList(db database.Database) []*entity.Province
}
