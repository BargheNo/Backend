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

func (addressService *AddressService) CreateAddress(addressInfo addressdto.CreateAddressRequest) (addressdto.AddressResponse, error) {
	province, err := addressService.addressRepository.GetProvinceByID(addressService.db, addressInfo.ProvinceID)
	if err != nil {
		return addressdto.AddressResponse{}, err
	}
	if province == nil {
		notFoundError := exception.NotFoundError{Item: addressService.constants.Field.Province}
		return addressdto.AddressResponse{}, &notFoundError
	}
	city, err := addressService.addressRepository.GetCityByID(addressService.db, addressInfo.CityID)
	if err != nil {
		return addressdto.AddressResponse{}, err
	}
	if city == nil {
		notFoundError := exception.NotFoundError{Item: addressService.constants.Field.City}
		return addressdto.AddressResponse{}, &notFoundError
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
	err = addressService.addressRepository.CreateAddress(addressService.db, address)
	if err != nil {
		return addressdto.AddressResponse{}, err
	}
	return addressdto.AddressResponse{
		ID:            address.ID,
		Province:      province.Name,
		City:          city.Name,
		StreetAddress: address.StreetAddress,
		PostalCode:    address.PostalCode,
		HouseNumber:   address.HouseNumber,
		Unit:          address.Unit,
	}, nil
}

func (addressService *AddressService) GetAddress(ownerID uint, ownerType string) (addressdto.AddressResponse, error) {
	address, err := addressService.addressRepository.GetOwnerAddress(addressService.db, ownerID, ownerType)
	if err != nil {
		return addressdto.AddressResponse{}, err
	}
	if address == nil {
		notFoundError := exception.NotFoundError{Item: addressService.constants.Field.Address}
		return addressdto.AddressResponse{}, &notFoundError
	}
	province, err := addressService.addressRepository.GetProvinceByID(addressService.db, address.ProvinceID)
	if err != nil {
		return addressdto.AddressResponse{}, err
	}
	city, err := addressService.addressRepository.GetCityByID(addressService.db, address.CityID)
	if err != nil {
		return addressdto.AddressResponse{}, err
	}
	response := addressdto.AddressResponse{
		ID:            address.ID,
		ProvinceID:    province.ID,
		Province:      province.Name,
		CityID:        city.ID,
		City:          city.Name,
		StreetAddress: address.StreetAddress,
		PostalCode:    address.PostalCode,
		HouseNumber:   address.HouseNumber,
		Unit:          address.Unit,
	}
	return response, nil
}

func (addressService *AddressService) GetAddresses(ownerAddressInfo addressdto.GetOwnerAddressesRequest) ([]addressdto.AddressResponse, error) {
	addressEntities, err := addressService.addressRepository.GetOwnerAddresses(addressService.db, ownerAddressInfo.OwnerID, ownerAddressInfo.OwnerType)
	if err != nil {
		return nil, err
	}
	addresses := make([]addressdto.AddressResponse, len(addressEntities))
	for i, address := range addressEntities {
		province, err := addressService.addressRepository.GetProvinceByID(addressService.db, address.ProvinceID)
		if err != nil {
			return nil, err
		}
		city, err := addressService.addressRepository.GetCityByID(addressService.db, address.CityID)
		if err != nil {
			return nil, err
		}
		addresses[i] = addressdto.AddressResponse{
			ID:            address.ID,
			Province:      province.Name,
			ProvinceID:    province.ID,
			City:          city.Name,
			CityID:        city.ID,
			StreetAddress: address.StreetAddress,
			PostalCode:    address.PostalCode,
			HouseNumber:   address.HouseNumber,
			Unit:          address.Unit,
		}
	}
	return addresses, nil
}

func (addressService *AddressService) DeleteAddress(addressID uint) error {
	address, err := addressService.addressRepository.GetAddressByID(addressService.db, addressID)
	if err != nil {
		return err
	}
	if address == nil {
		notFoundError := exception.NotFoundError{Item: addressService.constants.Field.Address}
		return &notFoundError
	}
	err = addressService.addressRepository.DeleteAddress(addressService.db, address)
	if err != nil {
		return err
	}
	return nil
}

func (addressService *AddressService) GetProvinceList() ([]addressdto.ProvinceResponse, error) {
	provinces, err := addressService.addressRepository.GetProvinceList(addressService.db)
	if err != nil {
		return nil, err
	}
	provincesList := make([]addressdto.ProvinceResponse, len(provinces))
	for i, province := range provinces {
		provincesList[i] = addressdto.ProvinceResponse{
			ID:   province.ID,
			Name: province.Name,
		}
	}
	return provincesList, nil
}

func (addressService *AddressService) GetCityProvinceCities(province addressdto.GetProvinceCitiesRequest) ([]addressdto.CityResponse, error) {
	cities, err := addressService.addressRepository.GetProvinceCities(addressService.db, province.ProvinceID)
	if err != nil {
		return nil, err
	}
	citiesList := make([]addressdto.CityResponse, len(cities))
	for i, city := range cities {
		citiesList[i] = addressdto.CityResponse{
			ID:   city.ID,
			Name: city.Name,
		}
	}
	return citiesList, nil
}
