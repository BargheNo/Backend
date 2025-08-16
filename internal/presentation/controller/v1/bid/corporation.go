package bid

import (
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	paymentdto "github.com/BargheNo/Backend/internal/application/dto/payment"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CorporationBidController struct {
	constants  *bootstrap.Constants
	pagination *bootstrap.Pagination
	bidService usecase.BidService
}

func NewCorporationBidController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	bidService usecase.BidService,
) *CorporationBidController {
	return &CorporationBidController{
		constants:  constants,
		pagination: pagination,
		bidService: bidService,
	}
}

func (bidController *CorporationBidController) GetBidStatuses(ctx *gin.Context) {
	statuses := bidController.bidService.GetCorporationBidStatuses()

	controller.Response(ctx, 200, "", statuses)
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

	var installmentPlanParams *paymentdto.InstallmentPlanRequest = nil
	if params.PaymentTerms.InstallmentPlan != nil {
		installmentPlanParams = &paymentdto.InstallmentPlanRequest{
			NumberOfMonths:    params.PaymentTerms.InstallmentPlan.NumberOfMonths,
			DownPaymentAmount: params.PaymentTerms.InstallmentPlan.DownPaymentAmount,
			MonthlyAmount:     params.PaymentTerms.InstallmentPlan.MonthlyAmount,
			Notes:             params.PaymentTerms.InstallmentPlan.Notes,
		}
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
	if err := bidController.bidService.SetBid(bidInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.setBid")
	controller.Response(ctx, 200, message, nil)
}

func (bidController *CorporationBidController) GetBids(ctx *gin.Context) {
	type getBidsParams struct {
		CorporationID uint   `uri:"corporationID" validate:"required"`
		Status        uint   `form:"status"`
		Query         string `form:"query"`
		Page          int    `form:"page"`
		PageSize      int    `form:"pageSize"`
		SortBy        uint   `form:"sortBy"`
		Asc           bool   `form:"asc"`
	}
	params := controller.Validated[getBidsParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, bidController.pagination.DefaultPage, bidController.pagination.DefaultPageSize)

	bidsRequest := biddto.GetCorporationBidsRequest{
		CorporationID: params.CorporationID,
		UserID:        userID.(uint),
		Status:        params.Status,
		Query:         params.Query,
		Offset:        offset,
		Limit:         limit,
		SortBy:        params.SortBy,
		Asc:           params.Asc,
	}
	bids, count, err := bidController.bidService.GetCorporationBids(bidsRequest)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(bids, count, offset, limit)

	controller.Response(ctx, 200, "", data)
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
	bid, err := bidController.bidService.GetCorporationBid(bidsRequest)
	if err != nil {
		panic(err)
	}

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
		PaymentMethod *uint `json:"method"`

		InstallmentPlan *installmentPlan `json:"installmentPlan" validate:"required_if=PaymentMethod 1"`
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

	var paymentTermsParams *paymentdto.UpdatePaymentTermsRequest = nil
	if params.PaymentTerms != nil {
		var installmentPlanParams *paymentdto.UpdateInstallmentPlanRequest = nil
		if params.PaymentTerms.InstallmentPlan != nil {
			installmentPlanParams = &paymentdto.UpdateInstallmentPlanRequest{
				NumberOfMonths:    params.PaymentTerms.InstallmentPlan.NumberOfMonths,
				DownPaymentAmount: params.PaymentTerms.InstallmentPlan.DownPaymentAmount,
				MonthlyAmount:     params.PaymentTerms.InstallmentPlan.MonthlyAmount,
				Notes:             params.PaymentTerms.InstallmentPlan.Notes,
			}
		}

		paymentTermsParams = &paymentdto.UpdatePaymentTermsRequest{
			PaymentMethod:   params.PaymentTerms.PaymentMethod,
			InstallmentPlan: installmentPlanParams,
		}
	}

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
		PaymentTerms:     paymentTermsParams,
	}
	if err := bidController.bidService.UpdateBid(updateBidInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateBid")
	controller.Response(ctx, 200, message, nil)
}

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
	if err := bidController.bidService.CancelBid(bidInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.cancelBid")
	controller.Response(ctx, 200, message, nil)
}
