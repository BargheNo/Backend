package service

import (
	reportdto "github.com/BargheNo/Backend/internal/application/dto/report"
)

type ReportService interface {
	CreateMaintenanceReport(requestInfo reportdto.CreateReportRequest)
}
