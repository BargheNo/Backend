package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type AddressRepository interface {
	CreateAddress(db database.Database, address *entity.Address) error
	GetAddressByID(db database.Database, id uint) (*entity.Address, bool)
	GetOwnerAddresses(db database.Database, ownerID uint, ownerType string) []*entity.Address
}
