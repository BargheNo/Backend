package service

import (
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	"github.com/BargheNo/Backend/internal/domain/entity"
)

type AddressService interface {
	CreateAddress(addressInfo addressdto.CreateAddressRequest) addressdto.AddressResponse
	GetAddress(addressID uint) *entity.Address
	GetAddresses(addressInfo addressdto.GetOwnerAddressesRequest) []addressdto.AddressResponse
}
