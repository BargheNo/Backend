package corporation

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/application/port"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralCorporationController struct {
	constants          *bootstrap.Constants
	corporationService port.CorporationService
}

func NewGeneralCorporationController(
	constants *bootstrap.Constants,
	corporationService port.CorporationService,
) *GeneralCorporationController {
	return &GeneralCorporationController{
		constants:          constants,
		corporationService: corporationService,
	}
}

func (corporationController *GeneralCorporationController) GetContactTypes(ctx *gin.Context) {
	contactTypes, err := corporationController.corporationService.GetContactTypes()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", contactTypes)
}

func (corporationController *GeneralCorporationController) GetCorporations(ctx *gin.Context) {
	corporations, err := corporationController.corporationService.GetAvailableCorporations()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", corporations)
}
