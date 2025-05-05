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

func (s *AddressServiceMock) CreateAddress(addressInfo addressdto.CreateAddressRequest) addressdto.AddressResponse {
	args := s.Called(addressInfo)
	return args.Get(0).(addressdto.AddressResponse)
}

func (s *AddressServiceMock) GetAddress(ownerID uint, ownerType string) addressdto.AddressResponse {
	args := s.Called(ownerID, ownerType)
	return args.Get(0).(addressdto.AddressResponse)
}

func (s *AddressServiceMock) GetAddresses(addressInfo addressdto.GetOwnerAddressesRequest) []addressdto.AddressResponse {
	args := s.Called(addressInfo)
	return args.Get(0).([]addressdto.AddressResponse)
}

func (s *AddressServiceMock) GetProvinceList() []addressdto.ProvinceResponse {
	args := s.Called()
	return args.Get(0).([]addressdto.ProvinceResponse)
}

func (s *AddressServiceMock) GetCityProvinceCities(province addressdto.GetProvinceCitiesRequest) []addressdto.CityResponse {
	args := s.Called(province)
	return args.Get(0).([]addressdto.CityResponse)
}

func (s *AddressServiceMock) DeleteAddress(addressID uint) {
	args := s.Called(addressID)
	if args.Get(0) != nil {
		panic(args.Get(0))
	}
}
