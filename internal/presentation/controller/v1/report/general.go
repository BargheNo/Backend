package report

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralReportController struct {
	constants     *bootstrap.Constants
	reportService usecase.ReportService
}

func NewGeneralReportController(
	constants *bootstrap.Constants,
	reportService usecase.ReportService,
) *GeneralReportController {
	return &GeneralReportController{
		constants:     constants,
		reportService: reportService,
	}
}

func (reportController *GeneralReportController) GetSortableFields(ctx *gin.Context) {
	columns := reportController.reportService.GetReportSortableColumns()
	controller.Response(ctx, 200, "", columns)
}
