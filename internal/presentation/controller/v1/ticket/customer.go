package ticket

import (
	"mime/multipart"
	"strconv"

	"github.com/BargheNo/Backend/bootstrap"
	ticketdto "github.com/BargheNo/Backend/internal/application/dto/ticket"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerTicketController struct {
	constants     *bootstrap.Constants
	ticketService service.TicketService
	pagination    *bootstrap.Pagination
}

func NewCustomerTicketController(
	constants *bootstrap.Constants,
	ticketService service.TicketService,
	pagination *bootstrap.Pagination,
) *CustomerTicketController {
	return &CustomerTicketController{
		constants:     constants,
		ticketService: ticketService,
		pagination:    pagination,
	}
}

func (ticketController *CustomerTicketController) CreateTicket(ctx *gin.Context) {
	type createTicketParams struct {
		Subject     string                `form:"subject" validate:"required"`
		Description string                `form:"description" validate:"required"`
		Image       *multipart.FileHeader `form:"image"`
	}

	params := controller.Validated[createTicketParams](ctx)
	subject, err := strconv.Atoi(params.Subject)
	if err != nil {
		subject = 2
	}
	userID, _ := ctx.Get(ticketController.constants.Context.ID)
	requestInfo := ticketdto.CreateTicketRequest{
		OwnerID:     userID.(uint),
		Subject:     enum.TicketSubject(subject),
		Description: params.Description,
		Image:       params.Image,
	}

	ticketController.ticketService.CreateTicket(requestInfo)

	trans := controller.GetTranslator(ctx, ticketController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createTicket")
	controller.Response(ctx, 200, message, nil)
}

func (ticketController *CustomerTicketController) GetTickets(ctx *gin.Context) {
	ownerID, _ := ctx.Get(ticketController.constants.Context.ID)
	params := controller.GetPagination(ctx, ticketController.pagination.DefaultPage, ticketController.pagination.DefaultPageSize)
	offset, limit := params.GetOffsetLimit()
	listInfo := ticketdto.TicketListRequest{
		OwnerID: ownerID.(uint),
		Offset:  offset,
		Limit:   limit,
	}

	tickets := ticketController.ticketService.GetCustomerTickets(listInfo)

	controller.Response(ctx, 200, "success", tickets)
}

func (ticketController *CustomerTicketController) GetComments(ctx *gin.Context) {
	type getCommentsParams struct {
		TicketID uint `uri:"ticketID" binding:"required"`
	}
	params := controller.Validated[getCommentsParams](ctx)
	pagination := controller.GetPagination(ctx, ticketController.pagination.DefaultPage, ticketController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()
	ownerID, _ := ctx.Get(ticketController.constants.Context.ID)
	listInfo := ticketdto.TicketCommentListRequest{
		TicketID: params.TicketID,
		OwnerID:  ownerID.(uint),
		Offset:   offset,
		Limit:    limit,
	}
	comments := ticketController.ticketService.GetTicketComments(listInfo)

	controller.Response(ctx, 200, "", comments)
}

func (ticketController *CustomerTicketController) CreateComment(ctx *gin.Context) {
	type createCommentParams struct {
		TicketID uint   `uri:"ticketID" validate:"required"`
		Body     string `json:"body" validate:"required"`
	}

	params := controller.Validated[createCommentParams](ctx)
	userID, _ := ctx.Get(ticketController.constants.Context.ID)
	requestInfo := ticketdto.CreateTicketCommentRequest{
		TicketID: params.TicketID,
		OwnerID:  userID.(uint),
		Body:     params.Body,
	}

	ticketController.ticketService.CreateTicketComment(requestInfo)

	trans := controller.GetTranslator(ctx, ticketController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createTicketComment")
	controller.Response(ctx, 200, message, nil)
}
