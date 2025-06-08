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

// read users from table name maybe ? or use enums or constants instead ?
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
		OwnerType:     addressController.constants.AddressOwners.User,
	}
	createdAddress, err := addressController.addressService.CreateAddress(addressRequestInfo)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, addressController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createAddress")
	controller.Response(ctx, 200, message, createdAddress)
}

func (addressController *CustomerAddressController) GetCustomerAddresses(ctx *gin.Context) {
	ownerID, _ := ctx.Get(addressController.constants.Context.ID)
	ownerInfo := addressdto.GetOwnerAddressesRequest{
		OwnerID:   ownerID.(uint),
		OwnerType: addressController.constants.AddressOwners.User,
	}
	addresses, err := addressController.addressService.GetAddresses(ownerInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", addresses)
}
