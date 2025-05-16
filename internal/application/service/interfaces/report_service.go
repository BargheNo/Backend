package service

import (
	reportdto "github.com/BargheNo/Backend/internal/application/dto/report"
)

type ReportService interface {
	CreateMaintenanceReport(requestInfo reportdto.CreateReportRequest)
	CreatePanelReport(requestInfo reportdto.CreateReportRequest)
	GetMaintenanceReport(reportID uint) reportdto.MaintenanceReportResponse
	GetPanelReport(reportID uint) reportdto.PanelReportResponse
	GetMaintenanceReports(requestInfo reportdto.ReportListRequest) []reportdto.MaintenanceReportResponse
	GetPanelReports(requestInfo reportdto.ReportListRequest) []reportdto.PanelReportResponse
	ResolveReport(requestInfo reportdto.ResolveReportRequest)
}
