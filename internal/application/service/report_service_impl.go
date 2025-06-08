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

func (reportService *ReportService) createReport(requestInfo reportdto.CreateReportRequest) (*entity.Report, error) {
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
		return nil, err
	}
	return report, nil
}

func (reportService *ReportService) CreateMaintenanceReport(requestInfo reportdto.CreateReportRequest) error {
	if err := reportService.maintenanceService.ValidateCustomerRecord(requestInfo.ObjectID, requestInfo.ReportedByID); err != nil {
		return err
	}

	report, err := reportService.createReport(requestInfo)
	if err != nil {
		return err
	}

	acceptedPermissions := []enum.PermissionType{enum.ReportViewAll, enum.PermissionAll}
	reportService.sendReportNotification(acceptedPermissions, report.ID, enum.MaintenanceReportCreated)
	return nil
}

func (reportService *ReportService) CreatePanelReport(requestInfo reportdto.CreateReportRequest) error {
	_, err := reportService.installationService.ValidatePanelOwnership(requestInfo.ObjectID, requestInfo.ReportedByID)
	if err != nil {
		return err
	}

	report, err := reportService.createReport(requestInfo)
	if err != nil {
		return err
	}

	acceptedPermissions := []enum.PermissionType{enum.ReportViewAll, enum.PermissionAll}
	reportService.sendReportNotification(acceptedPermissions, report.ID, enum.PanelReportCreated)
	return nil
}

func (reportService *ReportService) GetMaintenanceReport(reportID uint) (reportdto.MaintenanceReportResponse, error) {
	report, err := reportService.reportRepository.GetReportByID(reportService.db, reportID)
	if err != nil {
		return reportdto.MaintenanceReportResponse{}, err
	}
	if report == nil {
		return reportdto.MaintenanceReportResponse{}, exception.NotFoundError{Item: reportService.constants.Field.Report}
	}

	maintenanceRequest, err := reportService.maintenanceService.GetRequestByAdmin(report.ObjectID)
	if err != nil {
		return reportdto.MaintenanceReportResponse{}, err
	}

	return reportdto.MaintenanceReportResponse{
		ID:                 report.ID,
		Description:        report.Description,
		MaintenanceRequest: maintenanceRequest,
		Status:             report.Status.String(),
	}, nil
}

func (reportService *ReportService) GetPanelReport(reportID uint) (reportdto.PanelReportResponse, error) {
	report, err := reportService.reportRepository.GetReportByID(reportService.db, reportID)
	if err != nil {
		return reportdto.PanelReportResponse{}, err
	}
	if report == nil {
		return reportdto.PanelReportResponse{}, exception.NotFoundError{Item: reportService.constants.Field.Report}
	}

	panel, err := reportService.installationService.GetPanelByAdmin(report.ObjectID)
	if err != nil {
		return reportdto.PanelReportResponse{}, err
	}

	return reportdto.PanelReportResponse{
		ID:          report.ID,
		Panel:       panel,
		Description: report.Description,
		Status:      report.Status.String(),
	}, nil
}

func (reportService *ReportService) GetMaintenanceReports(requestInfo reportdto.ReportListRequest) ([]reportdto.MaintenanceReportResponse, error) {
	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	reports, err := reportService.reportRepository.GetReportsByObjectType(reportService.db, reportService.constants.ReportObjectTypes.Maintenance, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	reportResponses := make([]reportdto.MaintenanceReportResponse, len(reports))

	for i, report := range reports {
		maintenanceRequest, err := reportService.maintenanceService.GetRequestByAdmin(report.ObjectID)
		if err != nil {
			return nil, err
		}
		reportResponses[i] = reportdto.MaintenanceReportResponse{
			ID:                 report.ID,
			Description:        report.Description,
			MaintenanceRequest: maintenanceRequest,
			Status:             report.Status.String(),
		}
	}

	return reportResponses, nil
}

func (reportService *ReportService) GetPanelReports(requestInfo reportdto.ReportListRequest) ([]reportdto.PanelReportResponse, error) {
	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	reports, err := reportService.reportRepository.GetReportsByObjectType(reportService.db, reportService.constants.ReportObjectTypes.Panel, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	reportResponses := make([]reportdto.PanelReportResponse, len(reports))
	for i, report := range reports {
		panel, err := reportService.installationService.GetPanelByAdmin(report.ObjectID)
		if err != nil {
			return nil, err
		}
		reportResponses[i] = reportdto.PanelReportResponse{
			ID:          report.ID,
			Panel:       panel,
			Description: report.Description,
			Status:      report.Status.String(),
		}
	}

	return reportResponses, nil
}

func (reportService *ReportService) ResolveReport(requestInfo reportdto.ResolveReportRequest) error {
	reportService.userService.GetUserCredential(requestInfo.UserID)
	report, err := reportService.reportRepository.GetReportByID(reportService.db, requestInfo.ReportID)
	if err != nil {
		return err
	}
	if report == nil {
		return exception.NotFoundError{Item: reportService.constants.Field.Report}
	}
	if report.Status == enum.ReportStatusResolved {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(reportService.constants.Field.Report, reportService.constants.Tag.AlreadyResolved)
		return conflictErrors
	}
	report.Status = enum.ReportStatusResolved
	err = reportService.reportRepository.UpdateReport(reportService.db, report)
	if err != nil {
		return err
	}
	return nil
}
