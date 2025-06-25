package ticket

import (
	"github.com/BargheNo/Backend/bootstrap"
	ticketdto "github.com/BargheNo/Backend/internal/application/dto/ticket"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminTicketController struct {
	constant      *bootstrap.Constants
	pagination    *bootstrap.Pagination
	userService   usecase.UserService
	ticketService usecase.TicketService
}

func NewAdminTicketController(
	constant *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	userService usecase.UserService,
	ticketService usecase.TicketService,
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

	tickets, err := ticketController.ticketService.GetAdminTickets(requestInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", tickets)
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

	tickets, err := ticketController.ticketService.GetAdminTicketComments(requestInfo)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", tickets)
}

func (ticketController *AdminTicketController) CreateComment(ctx *gin.Context) {
	type CreateCommentRequest struct {
		TicketID uint   `uri:"ticketID" validate:"required"`
		Body     string `json:"body" validate:"required"`
	}
	params := controller.Validated[CreateCommentRequest](ctx)
	ownerID, _ := ctx.Get(ticketController.constant.Context.ID)
	requestInfo := ticketdto.CreateTicketCommentRequest{
		TicketID:  params.TicketID,
		OwnerID:   ownerID.(uint),
		OwnerType: ticketController.constant.TicketCommentOwners.Admin,
		Body:      params.Body,
	}
	if err := ticketController.ticketService.CreateAdminTicketComment(requestInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, ticketController.constant.Context.Translator)
	message, _ := trans.Translate("successMessage.createTicketComment")
	controller.Response(ctx, 200, message, nil)
}

func (ticketController *AdminTicketController) ResolveTicket(ctx *gin.Context) {
	type ResolveTicketRequest struct {
		TicketID uint `uri:"ticketID" validate:"required"`
	}
	params := controller.Validated[ResolveTicketRequest](ctx)
	ownerID, _ := ctx.Get(ticketController.constant.Context.ID)
	requestInfo := ticketdto.ResolveTicketRequest{
		TicketID: params.TicketID,
		OwnerID:  ownerID.(uint),
	}
	if err := ticketController.ticketService.ResolveTicket(requestInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, ticketController.constant.Context.Translator)
	message, _ := trans.Translate("successMessage.ticketResolved")
	controller.Response(ctx, 200, message, nil)
}
