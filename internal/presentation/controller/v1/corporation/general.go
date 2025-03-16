package corporation

import (
	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralCorporationController struct {
	constants          *bootstrap.Constants
	corporationService service.CorporationService
}

func NewGeneralCorporationController(
	constants *bootstrap.Constants,
	corporationService service.CorporationService,
) *GeneralCorporationController {
	return &GeneralCorporationController{
		constants:          constants,
		corporationService: corporationService,
	}
}

func (corporationController *GeneralCorporationController) Register(ctx *gin.Context) {
	type registerParams struct {
		CIN string `json:"cin" validate:"required"`
	}
	params := controller.Validated[registerParams](ctx, &corporationController.constants.Context)
	registerInfo := corporationdto.RegisterRequest{
		CIN: params.CIN,
	}
	corporationController.corporationService.Register(registerInfo)
}
