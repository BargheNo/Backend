package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	reportdto "github.com/BargheNo/Backend/internal/application/dto/report"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
)

type ReportService struct {
	constants          *bootstrap.Constants
	userService        service.UserService
	reportRepository   repository.ReportRepository
	maintenanceService service.MaintenanceService
	db                 database.Database
}

func NewReportService(
	constants *bootstrap.Constants,
	userService service.UserService,
	reportRepository repository.ReportRepository,
	maintenanceService service.MaintenanceService,
	db database.Database,
) *ReportService {
	return &ReportService{
		constants:          constants,
		userService:        userService,
		reportRepository:   reportRepository,
		maintenanceService: maintenanceService,
		db:                 db,
	}
}

func (reportService *ReportService) CreateMaintenanceReport(requestInfo reportdto.CreateReportRequest) {
	reportService.userService.GetUserCredential(requestInfo.ReportedByID)
	reportService.maintenanceService.GetMaintenanceRecordByID(requestInfo.ObjectID)
	report := &entity.Report{
		ObjectID:       requestInfo.ObjectID,
		ObjectType:     requestInfo.ObjectType,
		ReportedByID:   requestInfo.ReportedByID,
		ReportedByType: requestInfo.ReportedByType,
		Description:    requestInfo.Description,
		Status:         enum.ReportStatusPending,
	}

	err := reportService.reportRepository.CreateReport(reportService.db, report)
	if err != nil {
		panic(err)
	}
}

func (reportService *ReportService) GetAdminReports(requestInfo reportdto.ReportListRequest) []reportdto.MaintenanceReportResponse {
	reportService.userService.GetUserCredential(requestInfo.OwnerID)
	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	reports := reportService.reportRepository.GetReports(reportService.db, paginationModifier, sortingModifier)
	reportResponses := make([]reportdto.MaintenanceReportResponse, len(reports))
	for i, report := range reports {
		maintenanceRecord := reportService.maintenanceService.GetMaintenanceRecordByID(report.ObjectID)
		reportResponses[i] = reportdto.MaintenanceReportResponse{
			ID:                report.ID,
			Description:       report.Description,
			MaintenanceRecord: maintenanceRecord,
		}
	}

	return reportResponses
}
