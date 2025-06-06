package bid

import (
	"github.com/BargheNo/Backend/bootstrap"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminBidController struct {
	constants  *bootstrap.Constants
	pagination *bootstrap.Pagination
	BidService service.BidService
}

func NewAdminBidController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	BidService service.BidService,
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
	bids := bidController.BidService.GetRequestBidsByAdmin(bidsRequest)

	controller.Response(ctx, 200, "", bids)
}
