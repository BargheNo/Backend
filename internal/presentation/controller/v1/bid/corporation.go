package bid

import (
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	paymentdto "github.com/BargheNo/Backend/internal/application/dto/payment"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CorporationBidController struct {
	constants  *bootstrap.Constants
	pagination *bootstrap.Pagination
	bidService service.BidService
}

func NewCorporationBidController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	bidService service.BidService,
) *CorporationBidController {
	return &CorporationBidController{
		constants:  constants,
		pagination: pagination,
		bidService: bidService,
	}
}

func (bidController *CorporationBidController) SetBid(ctx *gin.Context) {
	type installmentPlan struct {
		NumberOfMonths    uint   `json:"numberOfMonths" validate:"required"`
		DownPaymentAmount uint   `json:"downPaymentAmount" validate:"required"`
		MonthlyAmount     uint   `json:"monthlyAmount" validate:"required"`
		Notes             string `json:"notes" validate:"required"`
	}
	type paymentTerms struct {
		PaymentMethod   uint             `json:"method" validate:"required"`
		InstallmentPlan *installmentPlan `json:"installmentPlan" validate:"required_if=PaymentMethod 1"`
	}
	type setBidParams struct {
		CorporationID    uint         `uri:"corporationID" validate:"required"`
		RequestID        uint         `uri:"requestID" validate:"required"`
		Cost             uint         `json:"cost" validate:"required"`
		Area             uint         `json:"area" validate:"required"`
		Power            uint         `json:"power" validate:"required"`
		Description      string       `json:"description"`
		InstallationTime time.Time    `json:"installationTime" validate:"required"`
		GuaranteeID      *uint        `json:"guaranteeID"`
		PaymentTerms     paymentTerms `json:"paymentTerms" validate:"required"`
	}
	params := controller.Validated[setBidParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)

	var installmentPlanParams *paymentdto.InstallmentPlanRequest
	if params.PaymentTerms.InstallmentPlan != nil {
		installmentPlanParams = &paymentdto.InstallmentPlanRequest{
			NumberOfMonths:    params.PaymentTerms.InstallmentPlan.NumberOfMonths,
			DownPaymentAmount: params.PaymentTerms.InstallmentPlan.DownPaymentAmount,
			MonthlyAmount:     params.PaymentTerms.InstallmentPlan.MonthlyAmount,
			Notes:             params.PaymentTerms.InstallmentPlan.Notes,
		}
	} else {
		installmentPlanParams = nil
	}

	bidInfo := biddto.SetBidRequest{
		CorporationID:    params.CorporationID,
		RequestID:        params.RequestID,
		BidderID:         userID.(uint),
		Status:           enum.BidStatusPending,
		Cost:             params.Cost,
		Area:             params.Area,
		Power:            params.Power,
		Description:      params.Description,
		InstallationTime: params.InstallationTime,
		GuaranteeID:      params.GuaranteeID,
		PaymentTerms: paymentdto.PaymentTermsRequest{
			PaymentMethod:   params.PaymentTerms.PaymentMethod,
			InstallmentPlan: installmentPlanParams,
		},
	}
	bidController.bidService.SetBid(bidInfo)

	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.setBid")
	controller.Response(ctx, 200, message, nil)
}

func (bidController *CorporationBidController) GetBids(ctx *gin.Context) {
	type getBidsParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		Status        uint `form:"status" validate:"required"`
	}
	params := controller.Validated[getBidsParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)
	pagination := controller.GetPagination(ctx, bidController.pagination.DefaultPage, bidController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	bidsRequest := biddto.GetCorporationBidsRequest{
		CorporationID: params.CorporationID,
		UserID:        userID.(uint),
		Status:        params.Status,
		Offset:        offset,
		Limit:         limit,
	}
	bids := bidController.bidService.GetCorporationBids(bidsRequest)

	controller.Response(ctx, 200, "", bids)
}

func (bidController *CorporationBidController) GetBid(ctx *gin.Context) {
	type getBidsParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		BidID         uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[getBidsParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)

	bidsRequest := biddto.GetBidRequest{
		CorporationID: params.CorporationID,
		UserID:        userID.(uint),
		BidID:         params.BidID,
	}
	bid := bidController.bidService.GetCorporationBid(bidsRequest)

	controller.Response(ctx, 200, "", bid)
}

func (bidController *CorporationBidController) UpdateBid(ctx *gin.Context) {
	type installmentPlan struct {
		NumberOfMonths    *uint   `json:"numberOfMonths"`
		DownPaymentAmount *uint   `json:"downPaymentAmount"`
		MonthlyAmount     *uint   `json:"monthlyAmount"`
		Notes             *string `json:"notes"`
	}

	type paymentTerms struct {
		PaymentMethod   *uint            `json:"method"`
		InstallmentPlan *installmentPlan `json:"installmentPlan" validate:"required_if=method 1"`
	}

	type updateBidParams struct {
		CorporationID    uint          `uri:"corporationID" validate:"required"`
		BidID            uint          `uri:"bidID" validate:"required"`
		Cost             *uint         `json:"cost"`
		Area             *uint         `json:"area"`
		Power            *uint         `json:"power"`
		Description      *string       `json:"description"`
		InstallationTime *time.Time    `json:"installationTime"`
		GuaranteeID      *uint         `json:"guaranteeID"`
		PaymentTerms     *paymentTerms `json:"paymentTerms"`
	}
	params := controller.Validated[updateBidParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)

	updateBidInfo := biddto.UpdateBidRequest{
		CorporationID:    params.CorporationID,
		BidID:            params.BidID,
		BidderID:         userID.(uint),
		Cost:             params.Cost,
		Area:             params.Area,
		Power:            params.Power,
		Description:      params.Description,
		InstallationTime: params.InstallationTime,
		GuaranteeID:      params.GuaranteeID,
		PaymentTerms: &paymentdto.UpdatePaymentTermsRequest{
			PaymentMethod: params.PaymentTerms.PaymentMethod,
			InstallmentPlan: &paymentdto.UpdateInstallmentPlanRequest{
				NumberOfMonths:    params.PaymentTerms.InstallmentPlan.NumberOfMonths,
				DownPaymentAmount: params.PaymentTerms.InstallmentPlan.DownPaymentAmount,
				MonthlyAmount:     params.PaymentTerms.InstallmentPlan.MonthlyAmount,
				Notes:             params.PaymentTerms.InstallmentPlan.Notes,
			},
		},
	}
	bidController.bidService.UpdateBid(updateBidInfo)

	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateBid")
	controller.Response(ctx, 200, message, nil)
}

// CHECK THIS AGAIN LATER
func (bidController *CorporationBidController) CancelBid(ctx *gin.Context) {
	type cancelBidParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		BidID         uint `uri:"bidID" validate:"required"`
	}
	params := controller.Validated[cancelBidParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)

	bidInfo := biddto.GetBidRequest{
		CorporationID: params.CorporationID,
		UserID:        userID.(uint),
		BidID:         params.BidID,
	}
	bidController.bidService.CancelBid(bidInfo)

	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.cancelBid")
	controller.Response(ctx, 200, message, nil)
}
