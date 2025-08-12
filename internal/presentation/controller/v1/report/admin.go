package report

import (
	"github.com/BargheNo/Backend/bootstrap"
	reportdto "github.com/BargheNo/Backend/internal/application/dto/report"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminReportController struct {
	constants     *bootstrap.Constants
	pagination    *bootstrap.Pagination
	reportService usecase.ReportService
}

func NewAdminReportController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	reportService usecase.ReportService,
) *AdminReportController {
	return &AdminReportController{
		constants:     constants,
		pagination:    pagination,
		reportService: reportService,
	}
}

func (reportController *AdminReportController) GetMaintenanceReports(ctx *gin.Context) {
	type GetMaintenanceReportsRequest struct {
		Status   uint `form:"status"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
		SortBy   uint `form:"sortBy"`
		Asc      bool `form:"asc"`
	}
	params := controller.Validated[GetMaintenanceReportsRequest](ctx)

	ownerID, _ := ctx.Get(reportController.constants.Context.ID)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, reportController.pagination.DefaultPage, reportController.pagination.DefaultPageSize)

	requestInfo := reportdto.ReportListRequest{
		OwnerID: ownerID.(uint),
		Status:  params.Status,
		Offset:  offset,
		Limit:   limit,
		SortBy:  params.SortBy,
		Asc:     params.Asc,
	}
	reports, count, err := reportController.reportService.GetMaintenanceReports(requestInfo)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(reports, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (reportController *AdminReportController) GetPanelReports(ctx *gin.Context) {
	type GetPanelReportsRequest struct {
		Status   uint `form:"status"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
		SortBy   uint `form:"sortBy"`
		Asc      bool `form:"asc"`
	}
	params := controller.Validated[GetPanelReportsRequest](ctx)

	ownerID, _ := ctx.Get(reportController.constants.Context.ID)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, reportController.pagination.DefaultPage, reportController.pagination.DefaultPageSize)

	requestInfo := reportdto.ReportListRequest{
		OwnerID: ownerID.(uint),
		Status:  params.Status,
		Offset:  offset,
		Limit:   limit,
		SortBy:  params.SortBy,
		Asc:     params.Asc,
	}
	reports, count, err := reportController.reportService.GetPanelReports(requestInfo)
	if err != nil {
		panic(err)
	}
	data := controller.NewPaginatedResponse(reports, count, offset, limit)

	controller.Response(ctx, 200, "", data)
}

func (reportController *AdminReportController) ResolveReport(ctx *gin.Context) {
	type ResolveReportRequest struct {
		ReportID uint `uri:"reportID" validate:"required"`
	}
	params := controller.Validated[ResolveReportRequest](ctx)
	userID, _ := ctx.Get(reportController.constants.Context.ID)

	requestInfo := reportdto.ResolveReportRequest{
		ReportID: params.ReportID,
		UserID:   userID.(uint),
	}
	if err := reportController.reportService.ResolveReport(requestInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, reportController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.reportResolved")
	controller.Response(ctx, 200, message, nil)
}

func (reportController *AdminReportController) SearchMaintenanceReports(ctx *gin.Context) {
	type SearchReportsRequest struct {
		Query    string `form:"query" validate:"required"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		SortBy   uint   `form:"sortBy"`
		Asc      bool   `form:"asc"`
	}
	params := controller.Validated[SearchReportsRequest](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, reportController.pagination.DefaultPage, reportController.pagination.DefaultPageSize)

	requestInfo := reportdto.SearchReportsRequest{
		Query:  params.Query,
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}

	reports, count, err := reportController.reportService.SearchMaintenanceReports(requestInfo)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(reports, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}

func (reportController *AdminReportController) SearchPanelReports(ctx *gin.Context) {
	type SearchReportsRequest struct {
		Query    string `form:"query" validate:"required"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
		SortBy   uint   `form:"sortBy"`
		Asc      bool   `form:"asc"`
	}

	params := controller.Validated[SearchReportsRequest](ctx)

	offset, limit := controller.GetOffsetLimit(params.Page, params.PageSize, reportController.pagination.DefaultPage, reportController.pagination.DefaultPageSize)

	requestInfo := reportdto.SearchReportsRequest{
		Query:  params.Query,
		Offset: offset,
		Limit:  limit,
		SortBy: params.SortBy,
		Asc:    params.Asc,
	}

	reports, count, err := reportController.reportService.SearchPanelReports(requestInfo)
	if err != nil {
		panic(err)
	}

	data := controller.NewPaginatedResponse(reports, count, offset, limit)
	controller.Response(ctx, 200, "", data)
}
