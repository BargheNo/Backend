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
	address := &entity.Address{
		Province:      addressInfo.Province,
		City:          addressInfo.City,
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
		Province:      address.Province,
		City:          address.City,
		StreetAddress: address.StreetAddress,
		PostalCode:    address.PostalCode,
		HouseNumber:   address.HouseNumber,
		Unit:          address.Unit,
	}
}

func (addressService *AddressService) GetAddress(addressID uint) addressdto.AddressResponse {
	address, exist := addressService.addressRepository.GetAddressByID(addressService.db, addressID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: addressService.constants.Field.Address}
		panic(notFoundError)
	}
	response := addressdto.AddressResponse{
		ID:            address.ID,
		Province:      address.Province,
		City:          address.City,
		StreetAddress: address.StreetAddress,
		PostalCode:    address.PostalCode,
		HouseNumber:   address.HouseNumber,
		Unit:          address.Unit,
	}
	return response
}

func (addressService *AddressService) GetAddresses(addressInfo addressdto.GetOwnerAddressesRequest) []addressdto.AddressResponse {
	addressEntities := addressService.addressRepository.GetOwnerAddresses(addressService.db, addressInfo.OwnerID, addressInfo.OwnerType)
	addresses := make([]addressdto.AddressResponse, len(addressEntities))
	for i, address := range addressEntities {
		addresses[i] = addressdto.AddressResponse{
			ID:            address.ID,
			Province:      address.Province,
			City:          address.City,
			StreetAddress: address.StreetAddress,
			PostalCode:    address.PostalCode,
			HouseNumber:   address.HouseNumber,
			Unit:          address.Unit,
		}
	}
	return addresses
}
