package usecase

import (
	reportdto "github.com/BargheNo/Backend/internal/application/dto/report"
)

type ReportService interface {
	GetReportSortableColumns() []reportdto.GetReportEnumResponse
	CreateMaintenanceReport(requestInfo reportdto.CreateReportRequest) error
	CreatePanelReport(requestInfo reportdto.CreateReportRequest) error
	GetMaintenanceReport(reportID uint) (reportdto.MaintenanceReportResponse, error)
	GetPanelReport(reportID uint) (reportdto.PanelReportResponse, error)
	GetMaintenanceReports(requestInfo reportdto.ReportListRequest) ([]reportdto.MaintenanceReportResponse, int64, error)
	GetPanelReports(requestInfo reportdto.ReportListRequest) ([]reportdto.PanelReportResponse, int64, error)
	ResolveReport(requestInfo reportdto.ResolveReportRequest) error
	GetReportStatuses() []reportdto.GetReportEnumResponse
	SearchMaintenanceReports(requestInfo reportdto.SearchReportsRequest) ([]reportdto.MaintenanceReportResponse, int64, error)
	SearchPanelReports(requestInfo reportdto.SearchReportsRequest) ([]reportdto.PanelReportResponse, int64, error)
}
