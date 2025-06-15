package bid

import (
	"github.com/BargheNo/Backend/bootstrap"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerBidController struct {
	constants  *bootstrap.Constants
	pagination *bootstrap.Pagination
	BidService service.BidService
}

func NewCustomerBidController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	BidService service.BidService,
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
	}
	params := controller.Validated[getBidsParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)
	pagination := controller.GetPagination(ctx, bidController.pagination.DefaultPage, bidController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	bidsRequest := biddto.GetListRequestBidsRequest{
		RequestID: params.RequestID,
		UserID:    userID.(uint),
		Offset:    offset,
		Limit:     limit,
	}
	bids, err := bidController.BidService.GetRequestAnonymousBids(bidsRequest)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", bids)
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
	err := bidController.BidService.AcceptBid(bidsRequest)
	if err != nil {
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
	err := bidController.BidService.RejectBid(bidsRequest)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.rejectBid")
	controller.Response(ctx, 201, message, nil)
}
