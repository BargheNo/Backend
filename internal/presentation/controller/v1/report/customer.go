package report

import (
	"github.com/BargheNo/Backend/bootstrap"
	reportdto "github.com/BargheNo/Backend/internal/application/dto/report"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerReportController struct {
	constants     *bootstrap.Constants
	reportService service.ReportService
}

func NewCustomerReportController(
	constants *bootstrap.Constants,
	reportService service.ReportService,
) *CustomerReportController {
	return &CustomerReportController{
		constants:     constants,
		reportService: reportService,
	}
}

func (reportController *CustomerReportController) CreateMaintenanceReport(ctx *gin.Context) {
	type createMaintenanceReportRequest struct {
		RecordID    uint   `uri:"recordID" validate:"required"`
		Description string `json:"description" validate:"required"`
	}
	params := controller.Validated[createMaintenanceReportRequest](ctx)
	userID, _ := ctx.Get(reportController.constants.Context.ID)
	requestInfo := reportdto.CreateReportRequest{
		ObjectID:       params.RecordID,
		ObjectType:     reportController.constants.ReportObjectTypes.Maintenance,
		ReportedByID:   userID.(uint),
		ReportedByType: reportController.constants.ReportOwners.User,
		Description:    params.Description,
	}

	reportController.reportService.CreateMaintenanceReport(requestInfo)

	trans := controller.GetTranslator(ctx, reportController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createReport")
	controller.Response(ctx, 200, message, nil)

}
