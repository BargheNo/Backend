package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type AddressService struct {
	constants         *bootstrap.Constants
	addressRepository repository.AddressRepository
	db                database.Database
}

func NewAddressService(
	constants *bootstrap.Constants,
	addressRepository repository.AddressRepository,
	db database.Database,
) *AddressService {
	return &AddressService{
		constants:         constants,
		addressRepository: addressRepository,
		db:                db,
	}
}

func (addressService *AddressService) CreateAddress(addressInfo addressdto.CreateAddressRequest) addressdto.AddressResponse {
	province, exist := addressService.addressRepository.GetProvinceByID(addressService.db, addressInfo.ProvinceID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: addressService.constants.Field.Province}
		panic(notFoundError)
	}
	city, exist := addressService.addressRepository.GetCityByID(addressService.db, addressInfo.CityID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: addressService.constants.Field.City}
		panic(notFoundError)
	}
	address := &entity.Address{
		ProvinceID:    addressInfo.ProvinceID,
		CityID:        addressInfo.CityID,
		StreetAddress: addressInfo.StreetAddress,
		PostalCode:    addressInfo.PostalCode,
		HouseNumber:   addressInfo.HouseNumber,
		Unit:          addressInfo.Unit,
		OwnerID:       addressInfo.OwnerID,
		OwnerType:     addressInfo.OwnerType,
	}
	err := addressService.addressRepository.CreateAddress(addressService.db, address)
	if err != nil {
		panic(err)
	}
	return addressdto.AddressResponse{
		ID:            address.ID,
		Province:      province.Name,
		City:          city.Name,
		StreetAddress: address.StreetAddress,
		PostalCode:    address.PostalCode,
		HouseNumber:   address.HouseNumber,
		Unit:          address.Unit,
	}
}

func (addressService *AddressService) GetAddress(ownerID uint, ownerType string) addressdto.AddressResponse {
	address, exist := addressService.addressRepository.GetOwnerAddress(addressService.db, ownerID, ownerType)
	if !exist {
		notFoundError := exception.NotFoundError{Item: addressService.constants.Field.Address}
		panic(notFoundError)
	}
	province, _ := addressService.addressRepository.GetProvinceByID(addressService.db, address.ProvinceID)
	city, _ := addressService.addressRepository.GetCityByID(addressService.db, address.CityID)
	response := addressdto.AddressResponse{
		ID:            address.ID,
		Province:      province.Name,
		City:          city.Name,
		StreetAddress: address.StreetAddress,
		PostalCode:    address.PostalCode,
		HouseNumber:   address.HouseNumber,
		Unit:          address.Unit,
	}
	return response
}

func (addressService *AddressService) GetAddresses(ownerAddressInfo addressdto.GetOwnerAddressesRequest) []addressdto.AddressResponse {
	addressEntities := addressService.addressRepository.GetOwnerAddresses(addressService.db, ownerAddressInfo.OwnerID, ownerAddressInfo.OwnerType)
	addresses := make([]addressdto.AddressResponse, len(addressEntities))
	for i, address := range addressEntities {
		province, _ := addressService.addressRepository.GetProvinceByID(addressService.db, address.ProvinceID)
		city, _ := addressService.addressRepository.GetCityByID(addressService.db, address.CityID)
		addresses[i] = addressdto.AddressResponse{
			ID:            address.ID,
			Province:      province.Name,
			City:          city.Name,
			StreetAddress: address.StreetAddress,
			PostalCode:    address.PostalCode,
			HouseNumber:   address.HouseNumber,
			Unit:          address.Unit,
		}
	}
	return addresses
}

func (addressService *AddressService) DeleteAddress(addressID uint) {
	address, exist := addressService.addressRepository.GetAddressByID(addressService.db, addressID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: addressService.constants.Field.Address}
		panic(notFoundError)
	}
	err := addressService.addressRepository.DeleteAddress(addressService.db, address)
	if err != nil {
		panic(err)
	}
}

func (addressService *AddressService) GetProvinceList() []addressdto.ProvinceResponse {
	provinces := addressService.addressRepository.GetProvinceList(addressService.db)
	provinceList := make([]addressdto.ProvinceResponse, len(provinces))
	for i, province := range provinces {
		provinceList[i] = addressdto.ProvinceResponse{
			ID:   province.ID,
			Name: province.Name,
		}
	}
	return provinceList
}

func (addressService *AddressService) GetCityProvinceCities(province addressdto.GetProvinceCitiesRequest) []addressdto.CityResponse {
	cities := addressService.addressRepository.GetProvinceCities(addressService.db, province.ProvinceID)
	citiesList := make([]addressdto.CityResponse, len(cities))
	for i, city := range cities {
		citiesList[i] = addressdto.CityResponse{
			ID:   city.ID,
			Name: city.Name,
		}
	}
	return citiesList
}
