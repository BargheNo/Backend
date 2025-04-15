package ticket

import (
	"github.com/BargheNo/Backend/bootstrap"
	ticketdto "github.com/BargheNo/Backend/internal/application/dto/ticket"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminTicketController struct {
	constant      *bootstrap.Constants
	pagination    *bootstrap.Pagination
	userService   service.UserService
	ticketService service.TicketService
}

func NewAdminTicketController(
	constant *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	userService service.UserService,
	ticketService service.TicketService,
) *AdminTicketController {
	return &AdminTicketController{
		constant:      constant,
		pagination:    pagination,
		userService:   userService,
		ticketService: ticketService,
	}
}

func (ticketController *AdminTicketController) GetTickets(ctx *gin.Context) {
	pagination := controller.GetPagination(ctx, ticketController.pagination.DefaultPage, ticketController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()
	ownerID, _ := ctx.Get(ticketController.constant.Context.ID)
	requestInfo := ticketdto.TicketListRequest{
		OwnerID: ownerID.(uint),
		Offset:  offset,
		Limit:   limit,
	}

	tickets := ticketController.ticketService.GetTickets(requestInfo)

	controller.Response(ctx, 200, "success", tickets)
}

func (ticketController *AdminTicketController) GetComments(ctx *gin.Context) {
	type GetCommentsRequest struct {
		TicketID uint `uri:"ticketID" validate:"required"`
	}
	params := controller.Validated[GetCommentsRequest](ctx)
	pagination := controller.GetPagination(ctx, ticketController.pagination.DefaultPage, ticketController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()
	ownerID, _ := ctx.Get(ticketController.constant.Context.ID)
	requestInfo := ticketdto.TicketCommentListRequest{
		TicketID: params.TicketID,
		OwnerID:  ownerID.(uint),
		Offset:   offset,
		Limit:    limit,
	}

	tickets := ticketController.ticketService.GetTicketComments(requestInfo)

	controller.Response(ctx, 200, "success", tickets)
}
