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
		requestID uint `uri:"requestID" validate:"required"`
	}
	params := controller.Validated[getBidsParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)

	bidsRequest := biddto.GetRequestBidsRequest{
		RequestID: params.requestID,
		UserID:    userID.(uint),
	}
	bids := bidController.BidService.GetRequestBids(bidsRequest)

	controller.Response(ctx, 200, "", bids)
}
