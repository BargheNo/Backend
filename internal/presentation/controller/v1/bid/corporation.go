package bid

import (
	"strconv"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type BidController struct {
	constants  *bootstrap.Constants
	pagination *bootstrap.Pagination
	BidService service.BidService
}

func NewBidController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	BidService service.BidService,
) *BidController {
	return &BidController{
		constants:  constants,
		pagination: pagination,
		BidService: BidService,
	}
}

func (bidController *BidController) SetBid(ctx *gin.Context) {
	type setBidParams struct {
		InstallationRequestID uint      `json:"installationRequestId" validate:"required"`
		Cost                  uint      `json:"cost" validate:"required"`
		Description           string    `json:"description"`
		InstallationDate      time.Time `json:"installationDate" validate:"required"`
	}
	params := controller.Validated[setBidParams](ctx)
	corporationID, _ := ctx.Get(bidController.constants.Context.ID)
	bidInfo := biddto.SetBidRequest{
		InstallationRequestID: params.InstallationRequestID,
		CorporationID:         corporationID.(uint),
		Cost:                  params.Cost,
		InstallationDate:      params.InstallationDate,
		Description:           params.Description,
	}
	bidController.BidService.SetBid(bidInfo)
	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.setBid")
	controller.Response(ctx, 200, message, nil)
}

func (bidController *BidController) CancelBid(ctx *gin.Context) {
	type cancelBidParams struct {
		BidID                 uint `json:"bidId" validate:"required"`
		InstallationRequestID uint `json:"installationRequestId" validate:"required"`
	}
	params := controller.Validated[cancelBidParams](ctx)
	corporationID, _ := ctx.Get(bidController.constants.Context.ID)
	bidInfo := biddto.CancelBidRequest{
		BidID:                 params.BidID,
		InstallationRequestID: params.InstallationRequestID,
		CorporationID:         corporationID.(uint),
	}
	bidController.BidService.CancelBid(bidInfo)

	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.cancelBid")
	controller.Response(ctx, 200, message, nil)
}

func (bidController *BidController) GetBids(ctx *gin.Context) {
	defaultPage, err := strconv.Atoi(bidController.pagination.DefaultPage)
	if err != nil {
		defaultPage = 1
	}
	defaultPageSize, err := strconv.Atoi(bidController.pagination.DefaultPageSize)
	if err != nil {
		defaultPageSize = 10
	}
	params := controller.GetPagination(ctx, defaultPage, defaultPageSize)
	offset, limit := params.GetOffsetLimit()
	corporationID, _ := ctx.Get(bidController.constants.Context.ID)
	bidsRequest := biddto.GetBidsRequest{
		CorporationID: corporationID.(uint),
		Offset:        offset,
		Limit:         limit,
	}
	bids := bidController.BidService.GetBids(bidsRequest)

	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.getBids")
	controller.Response(ctx, 200, message, bids)
}
