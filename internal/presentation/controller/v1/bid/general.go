package bid

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralBidController struct {
	constants  *bootstrap.Constants
	bidService usecase.BidService
}

func NewGeneralBidController(
	constants *bootstrap.Constants,
	bidService usecase.BidService,

) *GeneralBidController {
	return &GeneralBidController{
		constants:  constants,
		bidService: bidService,
	}
}

func (bidController *GeneralBidController) GetSortableFields(ctx *gin.Context) {
	columns := bidController.bidService.GetBidSortableColumns()
	controller.Response(ctx, 200, "", columns)
}
