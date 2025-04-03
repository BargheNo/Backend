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

type CorporationBidController struct {
	constants  *bootstrap.Constants
	pagination *bootstrap.Pagination
	BidService service.BidService
}

func NewCorporationBidController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	BidService service.BidService,
) *CorporationBidController {
	return &CorporationBidController{
		constants:  constants,
		pagination: pagination,
		BidService: BidService,
	}
}

func (bidController *CorporationBidController) SetBid(ctx *gin.Context) {
	type setBidParams struct {
		CorporationID         uint      `uri:"corporationID" validate:"required"`
		InstallationRequestID uint      `json:"installationRequestId" validate:"required"`
		Cost                  uint      `json:"cost" validate:"required"`
		Description           string    `json:"description"`
		InstallationTime      time.Time `json:"installationTime" validate:"required"`
	}
	params := controller.Validated[setBidParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)

	bidInfo := biddto.SetBidRequest{
		CorporationID:         params.CorporationID,
		BidderID:              userID.(uint),
		Cost:                  params.Cost,
		InstallationRequestID: params.InstallationRequestID,
		InstallationTime:      params.InstallationTime,
		Description:           params.Description,
	}
	bidController.BidService.SetBid(bidInfo)

	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.setBid")
	controller.Response(ctx, 200, message, nil)
}

func (bidController *CorporationBidController) CancelBid(ctx *gin.Context) {
	type cancelBidParams struct {
		CorporationID         uint `uri:"corporationID" validate:"required"`
		BidID                 uint `json:"bidId" validate:"required"`
		InstallationRequestID uint `json:"installationRequestId" validate:"required"`
	}
	params := controller.Validated[cancelBidParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)

	bidInfo := biddto.CancelBidRequest{
		CorporationID:         params.CorporationID,
		BidderID:              userID.(uint),
		BidID:                 params.BidID,
		InstallationRequestID: params.InstallationRequestID,
	}
	bidController.BidService.CancelBid(bidInfo)

	trans := controller.GetTranslator(ctx, bidController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.cancelBid")
	controller.Response(ctx, 200, message, nil)
}

func (bidController *CorporationBidController) GetBids(ctx *gin.Context) {
	type getBidsParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[getBidsParams](ctx)
	userID, _ := ctx.Get(bidController.constants.Context.ID)

	defaultPage, err := strconv.Atoi(bidController.pagination.DefaultPage)
	if err != nil {
		defaultPage = 1
	}
	defaultPageSize, err := strconv.Atoi(bidController.pagination.DefaultPageSize)
	if err != nil {
		defaultPageSize = 10
	}
	pagination := controller.GetPagination(ctx, defaultPage, defaultPageSize)
	offset, limit := pagination.GetOffsetLimit()

	bidsRequest := biddto.GetCorporationBidsRequest{
		CorporationID: params.CorporationID,
		UserID:        userID.(uint),
		Offset:        offset,
		Limit:         limit,
	}
	bids := bidController.BidService.GetCorporationBids(bidsRequest)

	controller.Response(ctx, 200, "", bids)
}
