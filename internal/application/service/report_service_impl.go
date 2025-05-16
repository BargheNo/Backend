package serviceimpl

import (
	"encoding/json"
	"log"

	"github.com/BargheNo/Backend/bootstrap"
	reportdto "github.com/BargheNo/Backend/internal/application/dto/report"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/message"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
)

type ReportService struct {
	constants           *bootstrap.Constants
	userService         service.UserService
	maintenanceService  service.MaintenanceService
	installationService service.InstallationService
	rabbitMQ            message.Broker
	reportRepository    repository.ReportRepository
	db                  database.Database
}

func NewReportService(
	constants *bootstrap.Constants,
	userService service.UserService,
	maintenanceService service.MaintenanceService,
	installationService service.InstallationService,
	rabbitMQ message.Broker,
	reportRepository repository.ReportRepository,
	db database.Database,
) *ReportService {
	return &ReportService{
		constants:           constants,
		userService:         userService,
		maintenanceService:  maintenanceService,
		installationService: installationService,
		rabbitMQ:            rabbitMQ,
		reportRepository:    reportRepository,
		db:                  db,
	}
}

func (reportService *ReportService) sendReportNotification(acceptedPermissions []enum.PermissionType, reportID uint, notificationType enum.NotificationType) {
	admins := reportService.userService.GetUsersByPermission(acceptedPermissions)
	for _, admin := range admins {
		additionalData := reportdto.ReportNotificationData{
			ReportID: reportID,
		}
		data, err := json.Marshal(additionalData)
		if err != nil {
			log.Println("Invalid data for message notification")
		}

		msg := struct {
			TypeName    enum.NotificationType `json:"typeName"`
			RecipientID uint                  `json:"recipientID"`
			Data        []byte                `json:"data"`
		}{
			TypeName:    notificationType,
			RecipientID: admin.ID,
			Data:        data,
		}

		if err := reportService.rabbitMQ.PublishMessage(reportService.constants.RabbitMQ.Events.SendNotification, msg); err != nil {
			log.Printf("error during send notification after bid: %v", err)
		}
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

	acceptedPermissions := []enum.PermissionType{enum.ReportViewAll, enum.PermissionAll}
	reportService.sendReportNotification(acceptedPermissions, report.ID, enum.MaintenanceReportCreated)
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

	acceptedPermissions := []enum.PermissionType{enum.ReportViewAll, enum.PermissionAll}
	reportService.sendReportNotification(acceptedPermissions, report.ID, enum.PanelReportCreated)
}

func (reportService *ReportService) GetMaintenanceReport(reportID uint) reportdto.MaintenanceReportResponse {
	report, exist := reportService.reportRepository.GetReportByID(reportService.db, reportID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: reportService.constants.Field.Report}
		panic(notFoundError)
	}
	maintenanceRecord := reportService.maintenanceService.GetMaintenanceRecordByID(report.ObjectID)
	reportResponse := reportdto.MaintenanceReportResponse{
		ID:                report.ID,
		Description:       report.Description,
		MaintenanceRecord: maintenanceRecord,
		Status:            report.Status.String(),
	}
	return reportResponse
}

func (reportService *ReportService) GetPanelReport(reportID uint) reportdto.PanelReportResponse {
	report, exist := reportService.reportRepository.GetReportByID(reportService.db, reportID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: reportService.constants.Field.Report}
		panic(notFoundError)
	}
	panel := reportService.installationService.GetPanelByID(report.ObjectID)
	reportResponse := reportdto.PanelReportResponse{
		ID:          report.ID,
		Panel:       panel,
		Description: report.Description,
		Status:      report.Status.String(),
	}
	return reportResponse
}

func (reportService *ReportService) GetMaintenanceReports(requestInfo reportdto.ReportListRequest) []reportdto.MaintenanceReportResponse {
	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	reports := reportService.reportRepository.GetReportsByObjectType(reportService.db, reportService.constants.ReportObjectTypes.Maintenance, paginationModifier, sortingModifier)
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

func (reportService *ReportService) GetPanelReports(requestInfo reportdto.ReportListRequest) []reportdto.PanelReportResponse {
	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	reports := reportService.reportRepository.GetReportsByObjectType(reportService.db, reportService.constants.ReportObjectTypes.Panel, paginationModifier, sortingModifier)
	reportResponses := make([]reportdto.PanelReportResponse, len(reports))
	for i, report := range reports {
		panel := reportService.installationService.GetPanelByID(report.ObjectID)
		reportResponses[i] = reportdto.PanelReportResponse{
			ID:          report.ID,
			Panel:       panel,
			Description: report.Description,
			Status:      report.Status.String(),
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
