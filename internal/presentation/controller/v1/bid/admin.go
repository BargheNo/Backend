package bid

import (
	"github.com/BargheNo/Backend/bootstrap"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminBidController struct {
	constants  *bootstrap.Constants
	pagination *bootstrap.Pagination
	BidService usecase.BidService
}

func NewAdminBidController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	BidService usecase.BidService,
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
		Page      int  `form:"page"`
		PageSize  int  `form:"pageSize"`
		SortBy    uint `form:"sortBy"`
		Asc       bool `form:"asc"`
	}
	params := controller.Validated[getBidsParams](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, bidController.pagination.DefaultPage, bidController.pagination.DefaultPageSize)

	bidsRequest := biddto.GetListRequestBidsRequestByAdmin{
		RequestID: params.RequestID,
		Offset:    offset,
		Limit:     limit,
		SortBy:    params.SortBy,
		Asc:       params.Asc,
	}
	bids, count, err := bidController.BidService.GetRequestBidsByAdmin(bidsRequest)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(bids, count, params.Page, params.PageSize)

	controller.Response(ctx, 200, "", data)
}
