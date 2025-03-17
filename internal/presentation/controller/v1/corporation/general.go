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
		Name     string `json:"name" validate:"required"`
		CIN      string `json:"cin" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	params := controller.Validated[registerParams](ctx, &corporationController.constants.Context)
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
	params := controller.Validated[loginParams](ctx, &corporationController.constants.Context)
	loginInfo := corporationdto.LoginRequest{
		CIN:      params.CIN,
		Password: params.Password,
	}
	corporationInfo := corporationController.corporationService.Login(loginInfo)
	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.Login")
	controller.Response(ctx, 200, message, corporationInfo)
}

func (corporationController *GeneralCorporationController) GetInstallationRequests(ctx *gin.Context) {
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	installation_requests := corporationController.corporationService.GetInstallationRequests(corporationID.(uint))

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.GetInstallationRequests")
	controller.Response(ctx, 200, message, installation_requests)
}

func (corporationController *GeneralCorporationController) SetBid(ctx *gin.Context) {
	type setBidParams struct {
		InstallationRequestID uint    `json:"installationRequest_id" validate:"required"`
		MinCost               float64 `json:"minCost" validate:"required"`
		MaxCost               float64 `json:"maxCost" validate:"required"`
		MinDeadline           string  `json:"minDeadline" validate:"required"`
		MaxDeadline           string  `json:"maxDeadline" validate:"required"`
		Description           string  `json:"description" validate:"required"`
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
	corporationController.corporationService.SetBid(bidInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.SetBid")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *GeneralCorporationController) CancelBid(ctx *gin.Context) {
	type cancelBidParams struct {
		BidID                 uint `json:"bidId" validate:"required"`
		InstallationRequestID uint `json:"installationRequestId" validate:"required"`
	}
	params := controller.Validated[cancelBidParams](ctx, &corporationController.constants.Context)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	bidInfo := corporationdto.CancelBidRequest{
		BidderID:              params.BidID,
		InstallationRequestID: params.InstallationRequestID,
		CorporationID:         corporationID.(uint),
	}
	corporationController.corporationService.CancelBid(bidInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.CancelBid")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *GeneralCorporationController) GetBids(ctx *gin.Context) {
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	bids := corporationController.corporationService.GetBids(corporationID.(uint))

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.GetBids")
	controller.Response(ctx, 200, message, bids)
}
