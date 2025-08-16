package ticket

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralTicketController struct {
	constants     *bootstrap.Constants
	ticketService usecase.TicketService
	pagination    *bootstrap.Pagination
}

func NewGeneralTicketController(
	constants *bootstrap.Constants,
	ticketService usecase.TicketService,
	pagination *bootstrap.Pagination,
) *GeneralTicketController {
	return &GeneralTicketController{
		constants:     constants,
		ticketService: ticketService,
		pagination:    pagination,
	}
}

func (ticketController *GeneralTicketController) GetTicketStatuses(ctx *gin.Context) {
	ticketStatuses := ticketController.ticketService.GetTicketStatuses()
	controller.Response(ctx, 200, "", ticketStatuses)
}

func (ticketController *GeneralTicketController) GetTicketSubjects(ctx *gin.Context) {
	ticketSubjects := ticketController.ticketService.GetTicketSubjects()
	controller.Response(ctx, 200, "", ticketSubjects)
}

func (ticketController *GeneralTicketController) GetSortableFields(ctx *gin.Context) {
	columns := ticketController.ticketService.GetTicketSortableColumns()
	controller.Response(ctx, 200, "", columns)
}
