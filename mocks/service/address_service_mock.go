package mocks

import (
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	"github.com/stretchr/testify/mock"
)

type AddressServiceMock struct {
	mock.Mock
}

func NewAddressServiceMock() *AddressServiceMock {
	return &AddressServiceMock{}
}

func (a *AddressServiceMock) CreateAddress(addressInfo addressdto.CreateAddressRequest) (addressdto.AddressResponse, error) {
	args := a.Called(addressInfo)
	return args.Get(0).(addressdto.AddressResponse), args.Error(1)
}

func (a *AddressServiceMock) GetAddress(ownerID uint, ownerType string) (addressdto.AddressResponse, error) {
	args := a.Called(ownerID, ownerType)
	return args.Get(0).(addressdto.AddressResponse), args.Error(1)
}

func (a *AddressServiceMock) GetAddresses(addressInfo addressdto.GetOwnerAddressesRequest) ([]addressdto.AddressResponse, error) {
	args := a.Called(addressInfo)
	return args.Get(0).([]addressdto.AddressResponse), args.Error(1)
}

func (a *AddressServiceMock) GetProvinceList() ([]addressdto.ProvinceResponse, error) {
	args := a.Called()
	return args.Get(0).([]addressdto.ProvinceResponse), args.Error(1)
}

func (a *AddressServiceMock) GetCityProvinceCities(province addressdto.GetProvinceCitiesRequest) ([]addressdto.CityResponse, error) {
	args := a.Called(province)
	return args.Get(0).([]addressdto.CityResponse), args.Error(1)
}

func (a *AddressServiceMock) DeleteAddress(addressID uint) error {
	args := a.Called(addressID)
	return args.Error(0)
}
