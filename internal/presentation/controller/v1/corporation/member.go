package corporation

import (
	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type MemberCorporationController struct {
	constants  *bootstrap.Constants
	corporationService service.CorporationService
	BidService service.BidService
}

func NewMemberCorporationController(
	constants *bootstrap.Constants,
	corporationService service.CorporationService,
	BidService service.BidService,
) *MemberCorporationController {
	return &MemberCorporationController{
		constants:  constants,
		corporationService: corporationService,
		BidService: BidService,
	}
}

func (corporationController *MemberCorporationController) AddContactInfo(ctx *gin.Context) {
	type addContactInfoParams struct {
		Phone string `json:"phone"`
		Email       string `json:"email"`
		Eitaa       string `json:"eitaa"`
		Bale        string `json:"bale"`
		Website     string `json:"website"`
		WhatsApp    string `json:"whatsApp"`
		Instagram   string `json:"instagram"`
		Telegram    string `json:"telegram"`
		Linkedin    string `json:"linkedin"`
	}

	params := controller.Validated[addContactInfoParams](ctx, &corporationController.constants.Context)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	contactInfo := corporationdto.ContactInfoRequest{
		Phone: params.Phone,
		Email:       params.Email,
		Eitaa:       params.Eitaa,
		Bale:        params.Bale,
		Website:     params.Website,
		WhatsApp:    params.WhatsApp,
		Instagram:   params.Instagram,
		Telegram:    params.Telegram,
		Linkedin:    params.Linkedin,
	}

	corporationController.corporationService.AddContactInfo(corporationID.(uint), contactInfo)

}

func (corporationController *MemberCorporationController) GetInstallationRequests(ctx *gin.Context) {
	pagination := controller.GetPagination(ctx, &corporationController.constants.Context)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	installation_requests := corporationController.BidService.GetInstallationRequests(corporationID.(uint), pagination.Page, pagination.PageSize)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.GetInstallationRequests")
	controller.Response(ctx, 200, message, installation_requests)
}

func (corporationController *MemberCorporationController) SetBid(ctx *gin.Context) {
	type setBidParams struct {
		InstallationRequestID uint    `json:"installationRequestId" validate:"required"`
		MinCost               float64 `json:"minCost" validate:"required"`
		MaxCost               float64 `json:"maxCost" validate:"required"`
		MinDeadline           string  `json:"minDeadline" validate:"required"`
		MaxDeadline           string  `json:"maxDeadline" validate:"required"`
		Description           string  `json:"description"`
		InstallationTime      string  `json:"installationTime" validate:"required"`
	}
	params := controller.Validated[setBidParams](ctx, &corporationController.constants.Context)

	bidInfo := corporationdto.SetBidRequest{
		InstallationRequestID: params.InstallationRequestID,
		MinCost:               params.MinCost,
		MaxCost:               params.MaxCost,
		MinDeadline:           params.MinDeadline,
		MaxDeadline:           params.MaxDeadline,
		Description:           params.Description,
		InstallationTime:      params.InstallationTime,
	}

	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	bidInfo.CorporationID = corporationID.(uint)
	corporationController.BidService.SetBid(bidInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.SetBid")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *MemberCorporationController) CancelBid(ctx *gin.Context) {
	type cancelBidParams struct {
		BidID                 uint `json:"bidId" validate:"required"`
		InstallationRequestID uint `json:"installationRequestId" validate:"required"`
	}
	params := controller.Validated[cancelBidParams](ctx, &corporationController.constants.Context)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	bidInfo := corporationdto.CancelBidRequest{
		BidID:                 params.BidID,
		InstallationRequestID: params.InstallationRequestID,
		CorporationID:         corporationID.(uint),
	}
	corporationController.BidService.CancelBid(bidInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.CancelBid")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *MemberCorporationController) GetBids(ctx *gin.Context) {
	pagination := controller.GetPagination(ctx, &corporationController.constants.Context)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	bids := corporationController.BidService.GetBids(corporationID.(uint), pagination.Page, pagination.PageSize)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.GetBids")
	controller.Response(ctx, 200, message, bids)
}
