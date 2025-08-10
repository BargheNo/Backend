package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type AddressRepository interface {
	CreateAddress(db database.Database, address *entity.Address) error
	CreateCity(db database.Database, city *entity.City) error
	CreateProvince(db database.Database, province *entity.Province) error
	DeleteAddress(db database.Database, address *entity.Address) error
	GetAddressByID(db database.Database, id uint) (*entity.Address, error)
	GetCityByID(db database.Database, id uint) (*entity.City, error)
	GetCityByName(db database.Database, name string) (*entity.City, error)
	GetOwnerAddress(db database.Database, ownerID uint, ownerType string) (*entity.Address, error)
	GetOwnerAddresses(db database.Database, ownerID uint, ownerType string) ([]*entity.Address, error)
	GetProvinceByID(db database.Database, id uint) (*entity.Province, error)
	GetProvinceByName(db database.Database, name string) (*entity.Province, error)
	GetProvinceCities(db database.Database, provinceID uint) ([]*entity.City, error)
	GetProvinceList(db database.Database) ([]*entity.Province, error)
}
