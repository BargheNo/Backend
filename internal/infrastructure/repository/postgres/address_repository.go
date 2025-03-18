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

func (repo *AddressRepository) CreateAddress(db database.Database, address *entity.Address) error {
	return db.GetDB().Create(&address).Error
}
