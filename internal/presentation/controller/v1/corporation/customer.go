package corporation

import (
	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerCorporationController struct {
	constants          *bootstrap.Constants
	pagination         *bootstrap.Pagination
	corporationService service.CorporationService
}

func NewCustomerCorporationController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	corporationService service.CorporationService,
) *CustomerCorporationController {
	return &CustomerCorporationController{
		constants:          constants,
		pagination:         pagination,
		corporationService: corporationService,
	}
}

func (corporationController *CustomerCorporationController) Register(ctx *gin.Context) {
	type signatory struct {
		Name               string `json:"name" validate:"required"`
		NationalCardNumber string `json:"nationalCardNumber" validate:"required"`
		Position           string `json:"position"`
	}
	type registerParams struct {
		Name               string      `json:"name" validate:"required"`
		RegistrationNumber string      `json:"registrationNumber" validate:"required"`
		NationalID         string      `json:"nationalID" validate:"required"`
		IBAN               string      `json:"iban"`
		Signatories        []signatory `json:"signatories" validate:"required"`
	}
	params := controller.Validated[registerParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)
	signatories := make([]corporationdto.Signatory, len(params.Signatories))
	for i, signatory := range params.Signatories {
		signatories[i] = corporationdto.Signatory{
			Name:               signatory.Name,
			NationalCardNumber: signatory.NationalCardNumber,
			Position:           signatory.Position,
		}
	}
	registerInfo := corporationdto.RegisterRequest{
		ApplicantID:        userID.(uint),
		Name:               params.Name,
		NationalID:         params.NationalID,
		RegistrationNumber: params.RegistrationNumber,
		IBAN:               params.IBAN,
		Signatories:        signatories,
	}

	corporationInfo := corporationController.corporationService.Register(registerInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.corporationRegister")
	controller.Response(ctx, 200, message, corporationInfo)
}
