package bid

import (
	"github.com/BargheNo/Backend/bootstrap"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerBidController struct {
	constants  *bootstrap.Constants
	pagination *bootstrap.Pagination
	BidService usecase.BidService
}

func NewCustomerBidController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	BidService usecase.BidService,
) *CustomerBidController {
	return &CustomerBidController{
		constants:  constants,
		pagination: pagination,
		BidService: BidService,
	}
}

func (bidController *CustomerBidController) GetBids(ctx *gin.Context) {
	type getBidsParams struct {
		RequestID uint `uri:"requestID" validate:"required"`
		Status    uint `form:"status"`
		Page      int  `form:"page"`
		PageSize  int  `form:"pageSize"`
		SortBy    uint `form:"sortBy"`
		Asc       bool `form:"asc"`
	}
	params := controller.Validated[getBidsParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, bidController.pagination.DefaultPage, bidController.pagination.DefaultPageSize)

	bidsRequest := biddto.GetListRequestBidsRequest{
		RequestID: params.RequestID,
		UserID:    userID.(uint),
		Status:    params.Status,
		Offset:    offset,
		Limit:     limit,
		SortBy:    params.SortBy,
		Asc:       params.Asc,
	}
	bids, count, err := bidController.BidService.GetRequestAnonymousBids(bidsRequest)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(bids, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (bidController *CustomerBidController) GetBid(ctx *gin.Context) {
	type getBidsParams struct {
		RequestID uint `uri:"requestID" validate:"required"`
		BidID     uint `uri:"bidID" validate:"required"`
	}
	params := controller.Validated[getBidsParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)

	bidsRequest := biddto.GetCustomerBidRequest{
		UserID:    userID.(uint),
		RequestID: params.RequestID,
		BidID:     params.BidID,
	}
	bids, err := bidController.BidService.GetRequestAnonymousBid(bidsRequest)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", bids)
}

func (bidController *CustomerBidController) AcceptBid(ctx *gin.Context) {
	type acceptBidsParams struct {
		RequestID uint `uri:"requestID" validate:"required"`
		BidID     uint `uri:"bidID" validate:"required"`
	}
	params := controller.Validated[acceptBidsParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)

	bidsRequest := biddto.GetCustomerBidRequest{
		RequestID: params.RequestID,
		BidID:     params.BidID,
		UserID:    userID.(uint),
	}
	if err := bidController.BidService.AcceptBid(bidsRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.acceptBid")
	controller.Response(ctx, 201, message, nil)
}

func (bidController *CustomerBidController) RejectBid(ctx *gin.Context) {
	type acceptBidsParams struct {
		RequestID uint `uri:"requestID" validate:"required"`
		BidID     uint `uri:"bidID" validate:"required"`
	}
	params := controller.Validated[acceptBidsParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)

	bidsRequest := biddto.GetCustomerBidRequest{
		RequestID: params.RequestID,
		BidID:     params.BidID,
		UserID:    userID.(uint),
	}
	if err := bidController.BidService.RejectBid(bidsRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.rejectBid")
	controller.Response(ctx, 201, message, nil)
}
