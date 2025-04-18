package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	reportdto "github.com/BargheNo/Backend/internal/application/dto/report"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
)

type ReportService struct {
	constants           *bootstrap.Constants
	userService         service.UserService
	reportRepository    repository.ReportRepository
	maintenanceService  service.MaintenanceService
	installationService service.InstallationService
	db                  database.Database
}

func NewReportService(
	constants *bootstrap.Constants,
	userService service.UserService,
	reportRepository repository.ReportRepository,
	maintenanceService service.MaintenanceService,
	installationService service.InstallationService,
	db database.Database,
) *ReportService {
	return &ReportService{
		constants:           constants,
		userService:         userService,
		reportRepository:    reportRepository,
		maintenanceService:  maintenanceService,
		installationService: installationService,
		db:                  db,
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

func (reportService *ReportService) CreatePanelReport(requestInfo reportdto.CreateReportRequest) {
	reportService.userService.GetUserCredential(requestInfo.ReportedByID)
	reportService.installationService.GetPanelByID(requestInfo.ObjectID)
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
			Status:            report.Status.String(),
		}
	}

	return reportResponses
}

func (reportService *ReportService) ResolveReport(requestInfo reportdto.ResolveReportRequest) {
	reportService.userService.GetUserCredential(requestInfo.UserID)
	report, exist := reportService.reportRepository.GetReportByID(reportService.db, requestInfo.ReportID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: reportService.constants.Field.Report}
		panic(notFoundError)
	}
	if report.Status == enum.ReportStatusResolved {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(reportService.constants.Field.Report, reportService.constants.Tag.AlreadyResolved)
		panic(conflictErrors)
	}
	report.Status = enum.ReportStatusResolved
	err := reportService.reportRepository.UpdateReport(reportService.db, report)
	if err != nil {
		panic(err)
	}
}
