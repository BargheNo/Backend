package service

import (
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
)

type AddressService interface {
	CreateAddress(addressInfo addressdto.CreateAddressRequest) addressdto.AddressResponse
	GetAddress(addressID uint) addressdto.AddressResponse
	GetAddresses(addressInfo addressdto.GetOwnerAddressesRequest) []addressdto.AddressResponse
	GetProvinceList() []addressdto.ProvinceResponse
	GetCityProvinceCities(province addressdto.GetProvinceCitiesRequest) []addressdto.CityResponse
}
