package bid

import (
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	paymentdto "github.com/BargheNo/Backend/internal/application/dto/payment"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminBidController struct {
	constants  *bootstrap.Constants
	pagination *bootstrap.Pagination
	bidService usecase.BidService
}

func NewAdminBidController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	bidService usecase.BidService,
) *AdminBidController {
	return &AdminBidController{
		constants:  constants,
		pagination: pagination,
		bidService: bidService,
	}
}

func (bidController *AdminBidController) GetRequestBids(ctx *gin.Context) {
	type getBidsParams struct {
		RequestID uint   `uri:"requestID" validate:"required"`
		Status    uint   `form:"status"`
		Query     string `form:"query"`
		Page      int    `form:"page"`
		PageSize  int    `form:"pageSize"`
		SortBy    uint   `form:"sortBy"`
		Asc       bool   `form:"asc"`
	}
	params := controller.Validated[getBidsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, bidController.pagination.DefaultPage, bidController.pagination.DefaultPageSize)
	bidsRequest := biddto.GetListRequestBidsRequestByAdmin{
		RequestID: params.RequestID,
		Status:    params.Status,
		Query:     params.Query,
		Offset:    offset,
		Limit:     limit,
		SortBy:    params.SortBy,
		Asc:       params.Asc,
	}
	bids, count, err := bidController.bidService.GetRequestBidsByAdmin(bidsRequest)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(bids, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (bidController *AdminBidController) GetBids(ctx *gin.Context) {
	type getBidsParams struct {
		Status   uint   `form:"status"`
		Query    string `form:"query"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		SortBy   uint   `form:"sortBy"`
		Asc      bool   `form:"asc"`
	}
	params := controller.Validated[getBidsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, bidController.pagination.DefaultPage, bidController.pagination.DefaultPageSize)
	bidsRequest := biddto.GetListBidsRequestByAdmin{
		Status: params.Status,
		Query:  params.Query,
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}
	bids, count, err := bidController.bidService.GetBidsByAdmin(bidsRequest)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(bids, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (bidController *AdminBidController) GetBid(ctx *gin.Context) {
	type getBidsParams struct {
		BidID uint `uri:"bidID" validate:"required"`
	}
	params := controller.Validated[getBidsParams](ctx)

	bid, err := bidController.bidService.GetBidByAdmin(params.BidID)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", bid)
}

func (bidController *AdminBidController) DeleteBid(ctx *gin.Context) {
	type getBidsParams struct {
		BidID uint `uri:"bidID" validate:"required"`
	}
	params := controller.Validated[getBidsParams](ctx)

	if err := bidController.bidService.DeleteBidByAdmin(params.BidID); err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "successMessage.deleteBid", nil)
}

func (bidController *AdminBidController) UpdateBid(ctx *gin.Context) {
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
		BidID:            params.BidID,
		Cost:             params.Cost,
		Area:             params.Area,
		Power:            params.Power,
		Description:      params.Description,
		InstallationTime: params.InstallationTime,
		GuaranteeID:      params.GuaranteeID,
		PaymentTerms:     paymentTermsParams,
	}
	if err := bidController.bidService.UpdateBidByAdmin(updateBidInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateBid")
	controller.Response(ctx, 200, message, nil)
}
