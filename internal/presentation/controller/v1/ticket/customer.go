package ticket

import (
	"mime/multipart"

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
}

func NewCustomerTicketController(
	constants *bootstrap.Constants,
	ticketService service.TicketService,
) *CustomerTicketController {
	return &CustomerTicketController{
		constants:     constants,
		ticketService: ticketService,
	}
}

func (ticketController *CustomerTicketController) CreateTicket(ctx *gin.Context) {
	type createTicketParams struct {
		Subject     uint                  `json:"subject" validate:"required"`
		Description string                `json:"description" validate:"required"`
		Image       *multipart.FileHeader `form:"image"`
	}

	params := controller.Validated[createTicketParams](ctx)
	userID, _ := ctx.Get(ticketController.constants.Context.ID)
	requestInfo := ticketdto.CreateTicketRequest{
		OwnerID:     userID.(uint),
		Subject:     enum.TicketSubject(params.Subject),
		Description: params.Description,
		Image:       params.Image,
	}

	ticketController.ticketService.CreateTicket(requestInfo)

	trans := controller.GetTranslator(ctx, ticketController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createTicket")
	controller.Response(ctx, 200, message, nil)
}
