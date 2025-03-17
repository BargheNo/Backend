package corporation

import (
	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type MemberCorporationController struct {
	constants          *bootstrap.Constants
	corporationService service.CorporationService
}

func NewMemberCorporationController(
	constants *bootstrap.Constants,
	corporationService service.CorporationService,
) *MemberCorporationController {
	return &MemberCorporationController{
		constants:          constants,
		corporationService: corporationService,
	}
}

func (corporationController *MemberCorporationController) GetInstallationRequests(ctx *gin.Context) {
	pagination := controller.GetPagination(ctx, &corporationController.constants.Context)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	installation_requests := corporationController.corporationService.GetInstallationRequests(corporationID.(uint), pagination.Page, pagination.PageSize)

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
	corporationController.corporationService.SetBid(bidInfo)

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
	corporationController.corporationService.CancelBid(bidInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.CancelBid")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *MemberCorporationController) GetBids(ctx *gin.Context) {
	pagination := controller.GetPagination(ctx, &corporationController.constants.Context)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	bids := corporationController.corporationService.GetBids(corporationID.(uint), pagination.Page, pagination.PageSize)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.GetBids")
	controller.Response(ctx, 200, message, bids)
}
