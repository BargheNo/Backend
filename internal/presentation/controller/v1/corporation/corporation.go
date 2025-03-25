package corporation

import (
	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CorporationController struct {
	constants          *bootstrap.Constants
	pagination         *bootstrap.Pagination
	corporationService service.CorporationService
}

func NewCorporationController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	corporationService service.CorporationService,
) *CorporationController {
	return &CorporationController{
		constants:          constants,
		pagination:         pagination,
		corporationService: corporationService,
	}
}

func (corporationController *CorporationController) GetCorporationInfo(ctx *gin.Context) {
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	corporationId := corporationdto.IDRequest{
		CorporationID: corporationID.(uint),
	}
	corporationInfo := corporationController.corporationService.GetCorporationInfo(corporationId)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.getCorporationInfo")
	controller.Response(ctx, 200, message, corporationInfo)
}

func (corporationController *CorporationController) ChangePassword(ctx *gin.Context) {
	type changePasswordParams struct {
		NewPassword     string `json:"newPassword" validate:"required"`
		ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=NewPassword"`
	}
	params := controller.Validated[changePasswordParams](ctx)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	changePasswordRequest := corporationdto.ChangePasswordRequest{
		CorporationID:   corporationID.(uint),
		NewPassword:     params.NewPassword,
		ConfirmPassword: params.ConfirmPassword,
	}
	corporationController.corporationService.ChangePassword(changePasswordRequest)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.changePassword")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CorporationController) UpdateContactInfo(ctx *gin.Context) {
	type addContactInfoParams struct {
		Phone     string `json:"phone"`
		Email     string `json:"email"`
		Eitaa     string `json:"eitaa"`
		Bale      string `json:"bale"`
		Website   string `json:"website"`
		WhatsApp  string `json:"whatsApp"`
		Instagram string `json:"instagram"`
		Telegram  string `json:"telegram"`
		Linkedin  string `json:"linkedin"`
	}

	params := controller.Validated[addContactInfoParams](ctx)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	contactInfo := corporationdto.ContactInfoRequest{
		CorporationID: corporationID.(uint),
		Phone:         params.Phone,
		Email:         params.Email,
		Eitaa:         params.Eitaa,
		Bale:          params.Bale,
		Website:       params.Website,
		WhatsApp:      params.WhatsApp,
		Instagram:     params.Instagram,
		Telegram:      params.Telegram,
		Linkedin:      params.Linkedin,
	}

	corporationController.corporationService.UpdateContactInfo(contactInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateContactInfo")
	controller.Response(ctx, 200, message, nil)

}
