package address

import (
	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralAddressController struct {
	constants      *bootstrap.Constants
	addressService service.AddressService
}

func NewGeneralAddressController(
	constants *bootstrap.Constants,
	addressService service.AddressService,
) *GeneralAddressController {
	return &GeneralAddressController{
		constants:      constants,
		addressService: addressService,
	}
}

func (addressController *GeneralAddressController) GetProvince(ctx *gin.Context) {
	provinceList := addressController.addressService.GetProvinceList()

	controller.Response(ctx, 200, "", provinceList)
}

func (addressController *GeneralAddressController) GetProvinceCities(ctx *gin.Context) {
	type getCitiesParams struct {
		ProvinceID uint `uri:"provinceID" validate:"required"`
	}
	params := controller.Validated[getCitiesParams](ctx)
	provinceInfo := addressdto.GetProvinceCitiesRequest{
		ProvinceID: params.ProvinceID,
	}
	citiesList := addressController.addressService.GetCityProvinceCities(provinceInfo)

	controller.Response(ctx, 200, "", citiesList)
}
