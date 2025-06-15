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

func (s *AddressServiceMock) CreateAddress(addressInfo addressdto.CreateAddressRequest) (addressdto.AddressResponse, error) {
	args := s.Called(addressInfo)
	return args.Get(0).(addressdto.AddressResponse), args.Error(1)
}

func (s *AddressServiceMock) GetAddress(ownerID uint, ownerType string) (addressdto.AddressResponse, error) {
	args := s.Called(ownerID, ownerType)
	return args.Get(0).(addressdto.AddressResponse), args.Error(1)
}

func (s *AddressServiceMock) GetAddresses(addressInfo addressdto.GetOwnerAddressesRequest) ([]addressdto.AddressResponse, error) {
	args := s.Called(addressInfo)
	return args.Get(0).([]addressdto.AddressResponse), args.Error(1)
}

func (s *AddressServiceMock) GetProvinceList() ([]addressdto.ProvinceResponse, error) {
	args := s.Called()
	return args.Get(0).([]addressdto.ProvinceResponse), args.Error(1)
}

func (s *AddressServiceMock) GetCityProvinceCities(province addressdto.GetProvinceCitiesRequest) ([]addressdto.CityResponse, error) {
	args := s.Called(province)
	return args.Get(0).([]addressdto.CityResponse), args.Error(1)
}

func (s *AddressServiceMock) DeleteAddress(addressID uint) error {
	args := s.Called(addressID)
	return args.Error(0)
}
