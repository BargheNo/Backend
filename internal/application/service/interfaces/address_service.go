package service

import (
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
)

type AddressService interface {
	CreateAddress(addressInfo addressdto.CreateAddressRequest) (addressdto.AddressResponse, error)
	GetAddress(ownerID uint, ownerType string) (addressdto.AddressResponse, error)
	GetAddresses(addressInfo addressdto.GetOwnerAddressesRequest) ([]addressdto.AddressResponse, error)
	GetProvinceList() ([]addressdto.ProvinceResponse, error)
	GetCityProvinceCities(province addressdto.GetProvinceCitiesRequest) ([]addressdto.CityResponse, error)
	DeleteAddress(addressID uint) error
}
