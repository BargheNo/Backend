package usecase

import (
	reportdto "github.com/BargheNo/Backend/internal/application/dto/report"
)

type ReportService interface {
	CreateMaintenanceReport(requestInfo reportdto.CreateReportRequest) error
	CreatePanelReport(requestInfo reportdto.CreateReportRequest) error
	GetMaintenanceReport(reportID uint) (reportdto.MaintenanceReportResponse, error)
	GetPanelReport(reportID uint) (reportdto.PanelReportResponse, error)
	GetMaintenanceReports(requestInfo reportdto.ReportListRequest) ([]reportdto.MaintenanceReportResponse, error)
	GetPanelReports(requestInfo reportdto.ReportListRequest) ([]reportdto.PanelReportResponse, error)
	ResolveReport(requestInfo reportdto.ResolveReportRequest) error
}
