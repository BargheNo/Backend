package service

import (
	"encoding/json"
	"log"

	"github.com/BargheNo/Backend/bootstrap"
	reportdto "github.com/BargheNo/Backend/internal/application/dto/report"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/enum/sortby"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/message"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type ReportService struct {
	constants           *bootstrap.Constants
	userService         usecase.UserService
	maintenanceService  usecase.MaintenanceService
	installationService usecase.InstallationService
	rabbitMQ            message.Broker
	reportRepository    postgres.ReportRepository
	db                  database.Database
}

func NewReportService(
	constants *bootstrap.Constants,
	userService usecase.UserService,
	maintenanceService usecase.MaintenanceService,
	installationService usecase.InstallationService,
	rabbitMQ message.Broker,
	reportRepository postgres.ReportRepository,
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

func (reportService *ReportService) getSortByColumn(requested uint) string {
	allowed := sortby.GetReportSortableColumns()
	_, ok := allowed[sortby.ReportSortBy(requested)]
	if !ok {
		return sortby.ReportSortByCreatedAt.DBColumn()
	}
	return sortby.ReportSortBy(requested).DBColumn()
}

func (reportService *ReportService) getReport(reportID uint) (*entity.Report, error) {
	report, err := reportService.reportRepository.FindReportByID(reportService.db, reportID)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, exception.NotFoundError{Item: reportService.constants.Field.Report}
	}
	return report, nil
}

func (reportService *ReportService) sendReportNotification(acceptedPermissions []enum.PermissionType, reportID uint, notificationType enum.NotificationType) {
	admins, err := reportService.userService.GetUsersByPermission(acceptedPermissions)
	if err != nil {
		log.Printf("error during send notification after bid: %v", err)
	}
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

func (reportService *ReportService) mapToFilterStatuses(enumStatus uint) []enum.ReportStatus {
	statuses := enum.GetAllReportStatuses()
	for _, status := range statuses {
		if uint(status) == enumStatus {
			if status == enum.ReportStatusAll {
				return statuses
			}
			return []enum.ReportStatus{status}
		}
	}
	return statuses
}

func (reportService *ReportService) GetReportSortableColumns() []reportdto.GetReportEnumResponse {
	columns := sortby.GetReportSortableColumns()
	response := make([]reportdto.GetReportEnumResponse, len(columns))
	i := 0
	for col, _ := range columns {
		response[i] = reportdto.GetReportEnumResponse{
			ID:   uint(col),
			Name: col.Name(),
		}
		i++
	}
	return response
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
	if _, err := reportService.installationService.ValidatePanelOwnership(requestInfo.ObjectID, requestInfo.ReportedByID); err != nil {
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
	report, err := reportService.getReport(reportID)
	if err != nil {
		return reportdto.MaintenanceReportResponse{}, err
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
	report, err := reportService.getReport(reportID)
	if err != nil {
		return reportdto.PanelReportResponse{}, err
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

func (reportService *ReportService) GetMaintenanceReports(requestInfo reportdto.ReportListRequest) ([]reportdto.MaintenanceReportResponse, int64, error) {
	options := postgres.NewQueryOptions().
		WithPagination(requestInfo.Limit, requestInfo.Offset).
		WithSorting(reportService.getSortByColumn(requestInfo.SortBy), requestInfo.Asc)

	statuses := reportService.mapToFilterStatuses(requestInfo.Status)

	reports, err := reportService.reportRepository.GetReportsByObjectType(reportService.db, reportService.constants.ReportObjectTypes.Maintenance, statuses, options)
	if err != nil {
		return nil, 0, err
	}
	reportResponses := make([]reportdto.MaintenanceReportResponse, len(reports))

	for i, report := range reports {
		maintenanceRequest, err := reportService.maintenanceService.GetRequestByAdmin(report.ObjectID)
		if err != nil {
			return nil, 0, err
		}

		reportResponses[i] = reportdto.MaintenanceReportResponse{
			ID:                 report.ID,
			Description:        report.Description,
			MaintenanceRequest: maintenanceRequest,
			Status:             report.Status.String(),
		}
	}

	count, err := reportService.reportRepository.CountReportsByObjectType(reportService.db, reportService.constants.ReportObjectTypes.Maintenance, statuses)
	if err != nil {
		return nil, 0, err
	}

	return reportResponses, count, nil
}

func (reportService *ReportService) GetPanelReports(requestInfo reportdto.ReportListRequest) ([]reportdto.PanelReportResponse, int64, error) {
	options := postgres.NewQueryOptions().
		WithPagination(requestInfo.Limit, requestInfo.Offset).
		WithSorting(reportService.getSortByColumn(requestInfo.SortBy), requestInfo.Asc)

	statuses := reportService.mapToFilterStatuses(requestInfo.Status)

	reports, err := reportService.reportRepository.GetReportsByObjectType(reportService.db, reportService.constants.ReportObjectTypes.Panel, statuses, options)
	if err != nil {
		return nil, 0, err
	}
	reportResponses := make([]reportdto.PanelReportResponse, len(reports))

	for i, report := range reports {
		panel, err := reportService.installationService.GetPanelByAdmin(report.ObjectID)
		if err != nil {
			return nil, 0, err
		}

		reportResponses[i] = reportdto.PanelReportResponse{
			ID:          report.ID,
			Panel:       panel,
			Description: report.Description,
			Status:      report.Status.String(),
		}
	}

	count, err := reportService.reportRepository.CountReportsByObjectType(reportService.db, reportService.constants.ReportObjectTypes.Maintenance, statuses)
	if err != nil {
		return nil, 0, err
	}

	return reportResponses, count, nil
}

func (reportService *ReportService) ResolveReport(requestInfo reportdto.ResolveReportRequest) error {
	report, err := reportService.getReport(requestInfo.ReportID)
	if err != nil {
		return err
	}

	if report.Status == enum.ReportStatusResolved {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(reportService.constants.Field.Report, reportService.constants.Tag.AlreadyResolved)
		return conflictErrors
	}

	report.Status = enum.ReportStatusResolved
	if err = reportService.reportRepository.UpdateReport(reportService.db, report); err != nil {
		return err
	}
	return nil
}
