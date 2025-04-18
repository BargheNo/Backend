package service

import (
	reportdto "github.com/BargheNo/Backend/internal/application/dto/report"
)

type ReportService interface {
	CreateMaintenanceReport(requestInfo reportdto.CreateReportRequest)
	CreatePanelReport(requestInfo reportdto.CreateReportRequest)
	GetAdminReports(requestInfo reportdto.ReportListRequest) []reportdto.MaintenanceReportResponse
	ResolveReport(requestInfo reportdto.ResolveReportRequest)
}
