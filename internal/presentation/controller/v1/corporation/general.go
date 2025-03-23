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
		Name          string `json:"name" validate:"required"`
		CIN           string `json:"cin" validate:"required"`
		Password      string `json:"password" validate:"required"`
		AcceptedTerms bool   `json:"acceptedTerms" validate:"eq=true"`
	}
	params := controller.Validated[registerParams](ctx)
	registerInfo := corporationdto.RegisterRequest{
		Name:     params.Name,
		CIN:      params.CIN,
		Password: params.Password,
	}
	corporationController.corporationService.Register(registerInfo)
	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.corporationRegister")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *GeneralCorporationController) Login(ctx *gin.Context) {
	type loginParams struct {
		CIN      string `json:"cin" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	params := controller.Validated[loginParams](ctx)
	loginInfo := corporationdto.LoginRequest{
		CIN:      params.CIN,
		Password: params.Password,
	}
	corporationInfo := corporationController.corporationService.Login(loginInfo)
	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.Login")
	controller.Response(ctx, 200, message, corporationInfo)
}
