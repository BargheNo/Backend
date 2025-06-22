package bid

import (
	"github.com/BargheNo/Backend/bootstrap"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	"github.com/BargheNo/Backend/internal/application/port"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminBidController struct {
	constants  *bootstrap.Constants
	pagination *bootstrap.Pagination
	BidService port.BidService
}

func NewAdminBidController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	BidService port.BidService,
) *AdminBidController {
	return &AdminBidController{
		constants:  constants,
		pagination: pagination,
		BidService: BidService,
	}
}

func (bidController *AdminBidController) GetBids(ctx *gin.Context) {
	type getBidsParams struct {
		RequestID uint `uri:"requestID" validate:"required"`
	}
	params := controller.Validated[getBidsParams](ctx)

	pagination := controller.GetPagination(ctx, bidController.pagination.DefaultPage, bidController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	bidsRequest := biddto.GetListRequestBidsRequestByAdmin{
		RequestID: params.RequestID,
		Offset:    offset,
		Limit:     limit,
	}
	bids, err := bidController.BidService.GetRequestBidsByAdmin(bidsRequest)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", bids)
}
