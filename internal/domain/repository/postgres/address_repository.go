package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type AddressRepository interface {
	CreateProvince(db database.Database, province *entity.Province) error
	CreateCity(db database.Database, city *entity.City) error
	CreateAddress(db database.Database, address *entity.Address) error
	GetProvinceList(db database.Database) []*entity.Province
	GetProvinceCities(db database.Database, provinceID uint) []*entity.City
	GetProvinceByID(db database.Database, id uint) (*entity.Province, bool)
	GetProvinceByName(db database.Database, name string) (*entity.Province, bool)
	GetCityByID(db database.Database, id uint) (*entity.City, bool)
	GetCityByName(db database.Database, name string) (*entity.City, bool)
	GetAddressByID(db database.Database, id uint) (*entity.Address, bool)
	GetOwnerAddresses(db database.Database, ownerID uint, ownerType string) []*entity.Address
	DeleteAddress(db database.Database, address *entity.Address) error
}
