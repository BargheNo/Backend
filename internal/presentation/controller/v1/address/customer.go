package address

import (
	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerAddressController struct {
	constants      *bootstrap.Constants
	addressService service.AddressService
}

func NewCustomerAddressController(
	constants *bootstrap.Constants,
	addressService service.AddressService,
) *CustomerAddressController {
	return &CustomerAddressController{
		constants:      constants,
		addressService: addressService,
	}
}

func (addressController *CustomerAddressController) CreateUserAddress(ctx *gin.Context) {
	type createAddressParams struct {
		ProvinceID    uint   `json:"provinceID" validate:"required"`
		CityID        uint   `json:"cityID" validate:"required"`
		StreetAddress string `json:"streetAddress" validate:"required"`
		PostalCode    string `json:"postalCode" validate:"required"`
		HouseNumber   string `json:"houseNumber" validate:"required"`
		Unit          uint   `json:"unit" validate:"required"`
	}
	params := controller.Validated[createAddressParams](ctx)
	ownerID, _ := ctx.Get(addressController.constants.Context.ID)
	addressRequestInfo := addressdto.CreateAddressRequest{
		ProvinceID:    params.ProvinceID,
		CityID:        params.CityID,
		StreetAddress: params.StreetAddress,
		PostalCode:    params.PostalCode,
		HouseNumber:   params.HouseNumber,
		Unit:          params.Unit,
		OwnerID:       ownerID.(uint),
		OwnerType:     "users",
	}
	createdAddress := addressController.addressService.CreateAddress(addressRequestInfo)

	trans := controller.GetTranslator(ctx, addressController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createAddress")
	controller.Response(ctx, 200, message, createdAddress)
}

func (addressController *CustomerAddressController) GetCustomerAddresses(ctx *gin.Context) {
	ownerID, _ := ctx.Get(addressController.constants.Context.ID)
	ownerInfo := addressdto.GetOwnerAddressesRequest{
		OwnerID:   ownerID.(uint),
		OwnerType: "users",
	}
	addresses := addressController.addressService.GetAddresses(ownerInfo)

	controller.Response(ctx, 200, "", addresses)
}

func (addressController *CustomerAddressController) CreateCorporationAddress(ctx *gin.Context) {
	type createAddressParams struct {
		ProvinceID    uint   `json:"provinceID" validate:"required"`
		CityID        uint   `json:"cityID" validate:"required"`
		StreetAddress string `json:"streetAddress" validate:"required"`
		PostalCode    string `json:"postalCode" validate:"required"`
		HouseNumber   string `json:"houseNumber" validate:"required"`
		Unit          uint   `json:"unit" validate:"required"`
	}
	params := controller.Validated[createAddressParams](ctx)
	ownerID, _ := ctx.Get(addressController.constants.Context.ID)
	addressRequestInfo := addressdto.CreateAddressRequest{
		ProvinceID:    params.ProvinceID,
		CityID:        params.CityID,
		StreetAddress: params.StreetAddress,
		PostalCode:    params.PostalCode,
		HouseNumber:   params.HouseNumber,
		Unit:          params.Unit,
		OwnerID:       ownerID.(uint),
		OwnerType:     "corporations",
	}
	createdAddress := addressController.addressService.CreateAddress(addressRequestInfo)

	trans := controller.GetTranslator(ctx, addressController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createAddress")
	controller.Response(ctx, 200, message, createdAddress)
}

func (addressController *CustomerAddressController) DeleteCorporationAddress(ctx *gin.Context) {
	type deleteAddressParams struct {
		AddressID uint `json:"addressID" validate:"required"`
	}
	params := controller.Validated[deleteAddressParams](ctx)
	addressController.addressService.DeleteAddress(params.AddressID)

	trans := controller.GetTranslator(ctx, addressController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteAddress")
	controller.Response(ctx, 200, message, nil)
}
