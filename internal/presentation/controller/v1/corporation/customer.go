package corporation

import (
	"strconv"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerCorporationController struct {
	constants          *bootstrap.Constants
	pagination         *bootstrap.Pagination
	corporationService service.CorporationService
	BidService         service.BidService
}

func NewCustomerCorporationController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	corporationService service.CorporationService,
	BidService service.BidService,
) *CustomerCorporationController {
	return &CustomerCorporationController{
		constants:          constants,
		pagination:         pagination,
		corporationService: corporationService,
		BidService:         BidService,
	}
}

func (corporationController *CustomerCorporationController) GetCorporationInfo(ctx *gin.Context) {
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	corporationId := corporationdto.IDRequest{
		CorporationID: corporationID.(uint),
	}
	corporationInfo := corporationController.corporationService.GetCorporationInfo(corporationId)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.getCorporationInfo")
	controller.Response(ctx, 200, message, corporationInfo)
}

func (corporationController *CustomerCorporationController) ChangePassword(ctx *gin.Context) {
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

func (corporationController *CustomerCorporationController) UpdateContactInfo(ctx *gin.Context) {
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

func (corporationController *CustomerCorporationController) SetBid(ctx *gin.Context) {
	type setBidParams struct {
		InstallationRequestID uint      `json:"installationRequestId" validate:"required"`
		Cost                  uint      `json:"cost" validate:"required"`
		Description           string    `json:"description"`
		InstallationDate      time.Time `json:"installationDate" validate:"required"`
	}
	params := controller.Validated[setBidParams](ctx)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	bidInfo := biddto.SetBidRequest{
		InstallationRequestID: params.InstallationRequestID,
		CorporationID:         corporationID.(uint),
		Cost:                  params.Cost,
		InstallationDate:      params.InstallationDate,
		Description:           params.Description,
	}
	corporationController.BidService.SetBid(bidInfo)
	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.setBid")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CustomerCorporationController) CancelBid(ctx *gin.Context) {
	type cancelBidParams struct {
		BidID                 uint `json:"bidId" validate:"required"`
		InstallationRequestID uint `json:"installationRequestId" validate:"required"`
	}
	params := controller.Validated[cancelBidParams](ctx)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	bidInfo := biddto.CancelBidRequest{
		BidID:                 params.BidID,
		InstallationRequestID: params.InstallationRequestID,
		CorporationID:         corporationID.(uint),
	}
	corporationController.BidService.CancelBid(bidInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.cancelBid")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CustomerCorporationController) GetBids(ctx *gin.Context) {
	defaultPage, err := strconv.Atoi(corporationController.pagination.DefaultPage)
	if err != nil {
		defaultPage = 1
	}
	defaultPageSize, err := strconv.Atoi(corporationController.pagination.DefaultPageSize)
	if err != nil {
		defaultPageSize = 10
	}
	params := controller.GetPagination(ctx, defaultPage, defaultPageSize)
	offset, limit := params.GetOffsetLimit()
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	bidsRequest := biddto.GetBidsRequest{
		CorporationID: corporationID.(uint),
		Offset:        offset,
		Limit:         limit,
	}
	bids := corporationController.BidService.GetBids(bidsRequest)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.getBids")
	controller.Response(ctx, 200, message, bids)
}
