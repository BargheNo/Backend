package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	maintenancedto "github.com/BargheNo/Backend/internal/application/dto/maintenance"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
)

type MaintenanceService struct {
	constants             *bootstrap.Constants
	userService           service.UserService
	installationService   service.InstallationService
	corporationService    service.CorporationService
	guaranteeService      service.GuaranteeService
	maintenanceRepository repository.MaintenanceRepository
	db                    database.Database
}

func NewMaintenanceService(
	constants *bootstrap.Constants,
	userService service.UserService,
	installationService service.InstallationService,
	corporationService service.CorporationService,
	guaranteeService service.GuaranteeService,
	maintenanceRepository repository.MaintenanceRepository,
	db database.Database,
) *MaintenanceService {
	return &MaintenanceService{
		constants:             constants,
		userService:           userService,
		installationService:   installationService,
		corporationService:    corporationService,
		guaranteeService:      guaranteeService,
		maintenanceRepository: maintenanceRepository,
		db:                    db,
	}
}

func (maintenanceService *MaintenanceService) ValidateCustomerRecord(recordID, userID uint) error {
	maintenanceRecord, err := maintenanceService.maintenanceRepository.FindRecordByID(maintenanceService.db, recordID)
	if err != nil {
		return err
	}
	if maintenanceRecord == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRecord}
		return notFoundError
	}

	maintenanceRequest, err := maintenanceService.maintenanceRepository.FindRequestByID(maintenanceService.db, maintenanceRecord.RequestID)
	if err != nil {
		return err
	}
	if maintenanceRequest == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		return notFoundError
	}

	_, err = maintenanceService.installationService.ValidatePanelOwnership(maintenanceRequest.PanelID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (maintenanceService *MaintenanceService) GetRequestByAdmin(recordID uint) (maintenancedto.AdminMaintenanceRequestResponse, error) {
	var response maintenancedto.AdminMaintenanceRequestResponse
	maintenanceRecord, err := maintenanceService.maintenanceRepository.FindRecordByID(maintenanceService.db, recordID)
	if err != nil {
		return response, err
	}
	if maintenanceRecord == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRecord}
		return response, notFoundError
	}

	maintenanceRequest, err := maintenanceService.maintenanceRepository.FindRequestByID(maintenanceService.db, maintenanceRecord.RequestID)
	if err != nil {
		return response, err
	}
	if maintenanceRequest == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		return response, notFoundError
	}

	panel, err := maintenanceService.installationService.GetPanelByAdmin(maintenanceRequest.PanelID)
	if err != nil {
		return response, err
	}

	corporation, err := maintenanceService.corporationService.GetCorporationCredentials(maintenanceRequest.CorporationID)
	if err != nil {
		return response, err
	}

	record, err := maintenanceService.getCorporationMaintenanceRecord(maintenanceRequest.ID, maintenanceRequest.PanelID)
	if err != nil {
		return response, err
	}

	response = maintenancedto.AdminMaintenanceRequestResponse{
		ID:                   maintenanceRequest.ID,
		CreatedAt:            maintenanceRequest.CreatedAt,
		Panel:                panel,
		Corporation:          corporation,
		Subject:              maintenanceRequest.Subject,
		Description:          maintenanceRequest.Description,
		UrgencyLevel:         maintenanceRequest.UrgencyLevel.String(),
		Status:               maintenanceRequest.Status.String(),
		IsGuaranteeRequested: maintenanceRequest.IsGuaranteeRequested,
		Record:               record,
	}

	return response, nil
}

func (maintenanceService *MaintenanceService) mapStatusForRole(statusID uint, agent enum.AgentType) []enum.MaintenanceRequestStatus {
	status := enum.MaintenanceRequestStatus(statusID)

	allowedStatuses := enum.GetAllowedMaintenanceRequestStatuses(agent)

	for _, allowedStatus := range allowedStatuses {
		if status == allowedStatus {
			if status == enum.MaintenanceRequestStatusAll {
				return allowedStatuses
			}
			return []enum.MaintenanceRequestStatus{status}
		}
	}
	return allowedStatuses
}

func (maintenanceService *MaintenanceService) GetMaintenanceUrgencyLevels() []maintenancedto.MaintenanceStatusesResponse {
	levels := enum.GetAllUrgencyLevels()
	response := make([]maintenancedto.MaintenanceStatusesResponse, len(levels))
	for i, level := range levels {
		response[i] = maintenancedto.MaintenanceStatusesResponse{
			ID:   uint(level),
			Name: level.String(),
		}
	}
	return response
}

func (maintenanceService *MaintenanceService) GetMaintenanceRequestStatuses(agent enum.AgentType) []maintenancedto.MaintenanceStatusesResponse {
	statuses := enum.GetAllowedMaintenanceRequestStatuses(agent)
	response := make([]maintenancedto.MaintenanceStatusesResponse, len(statuses))
	for i, status := range statuses {
		response[i] = maintenancedto.MaintenanceStatusesResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response
}

func (maintenanceService *MaintenanceService) CreateMaintenanceRequest(request maintenancedto.CreateMaintenanceRequest) error {
	if exist, err := maintenanceService.userService.IsUserActive(request.OwnerID); !exist {
		if err != nil {
			return err
		}
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: maintenanceService.constants.Field.MaintenanceRequest,
		}
		return forbiddenError
	}

	err := maintenanceService.corporationService.DoesCorporationExist(request.CorporationID)
	if err != nil {
		return err
	}

	_, err = maintenanceService.installationService.ValidatePanelOwnership(request.PanelID, request.OwnerID)
	if err != nil {
		return err
	}

	allowedStatus := []enum.MaintenanceRequestStatus{enum.MaintenanceRequestStatusPending}
	currentActiveRequest, err := maintenanceService.maintenanceRepository.FindRequestsByPanelIDAndStatus(maintenanceService.db, request.PanelID, allowedStatus)
	if err != nil {
		return err
	}
	if currentActiveRequest != nil {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.Pending)
		return conflictErrors
	}

	if request.IsUsingGuarantee {
		if err := maintenanceService.installationService.ValidatePanelGuarantee(request.PanelID); err != nil {
			return err
		}
	}

	maintenanceRequest := &entity.MaintenanceRequest{
		CorporationID:        request.CorporationID,
		PanelID:              request.PanelID,
		Subject:              request.Subject,
		Description:          request.Description,
		Status:               enum.MaintenanceRequestStatusPending,
		UrgencyLevel:         request.UrgencyLevel,
		IsGuaranteeRequested: request.IsUsingGuarantee,
	}
	err = maintenanceService.maintenanceRepository.CreateMaintenanceRequest(maintenanceService.db, maintenanceRequest)
	if err != nil {
		return err
	}
	return nil
}

func (maintenanceService *MaintenanceService) GetCustomerMaintenanceRequests(listInfo maintenancedto.CustomerMaintenanceListRequest) ([]maintenancedto.CustomerMaintenanceRequestResponse, error) {
	allowedStatus := maintenanceService.mapStatusForRole(listInfo.Status, enum.AgentTypeCustomer)

	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	maintenanceRequests, err := maintenanceService.maintenanceRepository.FindRequestsByCustomerID(maintenanceService.db, listInfo.OwnerID, allowedStatus, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	response := make([]maintenancedto.CustomerMaintenanceRequestResponse, len(maintenanceRequests))

	for i, maintenanceRequest := range maintenanceRequests {
		panelInfoRequest := installationdto.GetOwnerRequest{
			OwnerID:        listInfo.OwnerID,
			InstallationID: maintenanceRequest.PanelID,
		}
		panel, err := maintenanceService.installationService.GetCustomerPanel(panelInfoRequest)
		if err != nil {
			return nil, err
		}

		corporation, err := maintenanceService.corporationService.GetCorporationCredentials(maintenanceRequest.CorporationID)
		if err != nil {
			return nil, err
		}

		response[i] = maintenancedto.CustomerMaintenanceRequestResponse{
			ID:                   maintenanceRequest.ID,
			CreatedAt:            maintenanceRequest.CreatedAt,
			Panel:                panel,
			Corporation:          corporation,
			Subject:              maintenanceRequest.Subject,
			Description:          maintenanceRequest.Description,
			UrgencyLevel:         maintenanceRequest.UrgencyLevel.String(),
			Status:               maintenanceRequest.Status.String(),
			IsGuaranteeRequested: maintenanceRequest.IsGuaranteeRequested,
		}
	}
	return response, nil
}

func (maintenanceService *MaintenanceService) GetCustomerPanelMaintenanceRequests(listInfo maintenancedto.CustomerPanelMaintenanceListRequest) ([]maintenancedto.CustomerMaintenanceRequestResponse, error) {
	_, err := maintenanceService.installationService.ValidatePanelOwnership(listInfo.PanelID, listInfo.OwnerID)
	if err != nil {
		return nil, err
	}

	allowedStatus := maintenanceService.mapStatusForRole(listInfo.Status, enum.AgentTypeCustomer)

	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	maintenanceRequests, err := maintenanceService.maintenanceRepository.FindRequestsByPanelIDAndStatus(maintenanceService.db, listInfo.PanelID, allowedStatus, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	response := make([]maintenancedto.CustomerMaintenanceRequestResponse, len(maintenanceRequests))

	for i, maintenanceRequest := range maintenanceRequests {
		panelInfoRequest := installationdto.GetOwnerRequest{
			OwnerID:        listInfo.OwnerID,
			InstallationID: maintenanceRequest.PanelID,
		}
		panel, err := maintenanceService.installationService.GetCustomerPanel(panelInfoRequest)
		if err != nil {
			return nil, err
		}

		corporation, err := maintenanceService.corporationService.GetCorporationCredentials(maintenanceRequest.CorporationID)
		if err != nil {
			return nil, err
		}

		response[i] = maintenancedto.CustomerMaintenanceRequestResponse{
			ID:                   maintenanceRequest.ID,
			CreatedAt:            maintenanceRequest.CreatedAt,
			Panel:                panel,
			Corporation:          corporation,
			Subject:              maintenanceRequest.Subject,
			Description:          maintenanceRequest.Description,
			UrgencyLevel:         maintenanceRequest.UrgencyLevel.String(),
			Status:               maintenanceRequest.Status.String(),
			IsGuaranteeRequested: maintenanceRequest.IsGuaranteeRequested,
		}
	}
	return response, nil
}

func (maintenanceService *MaintenanceService) GetCustomerMaintenanceRequest(maintenanceInfo maintenancedto.CustomerMaintenanceRequest) (maintenancedto.CustomerMaintenanceRequestResponse, error) {
	maintenanceRequest, err := maintenanceService.maintenanceRepository.FindRequestByID(maintenanceService.db, maintenanceInfo.RequestID)
	if err != nil {
		return maintenancedto.CustomerMaintenanceRequestResponse{}, err
	}
	if maintenanceRequest == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		return maintenancedto.CustomerMaintenanceRequestResponse{}, notFoundError
	}

	panelInfoRequest := installationdto.GetOwnerRequest{
		OwnerID:        maintenanceInfo.OwnerID,
		InstallationID: maintenanceRequest.PanelID,
	}
	panel, err := maintenanceService.installationService.GetCustomerPanel(panelInfoRequest)
	if err != nil {
		return maintenancedto.CustomerMaintenanceRequestResponse{}, err
	}

	corporation, err := maintenanceService.corporationService.GetCorporationCredentials(maintenanceRequest.CorporationID)
	if err != nil {
		return maintenancedto.CustomerMaintenanceRequestResponse{}, err
	}

	record, err := maintenanceService.getCustomerMaintenanceRecord(maintenanceInfo.RequestID, maintenanceRequest.PanelID)
	if err != nil {
		return maintenancedto.CustomerMaintenanceRequestResponse{}, err
	}

	response := maintenancedto.CustomerMaintenanceRequestResponse{
		ID:                   maintenanceRequest.ID,
		CreatedAt:            maintenanceRequest.CreatedAt,
		Panel:                panel,
		Corporation:          corporation,
		Subject:              maintenanceRequest.Subject,
		Description:          maintenanceRequest.Description,
		UrgencyLevel:         maintenanceRequest.UrgencyLevel.String(),
		Status:               maintenanceRequest.Status.String(),
		IsGuaranteeRequested: maintenanceRequest.IsGuaranteeRequested,
		Record:               record,
	}
	return response, nil
}

func (maintenanceService *MaintenanceService) getCustomerMaintenanceRecord(requestID, panelID uint) (maintenancedto.CustomerMaintenanceRecordResponse, error) {
	record, err := maintenanceService.maintenanceRepository.FindRecordByRequestID(maintenanceService.db, requestID)
	if err != nil {
		return maintenancedto.CustomerMaintenanceRecordResponse{}, err
	}
	if record == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRecord}
		return maintenancedto.CustomerMaintenanceRecordResponse{}, notFoundError
	}

	violation, err := maintenanceService.guaranteeService.GetCustomerPanelGuaranteeViolation(panelID)
	if err != nil {
		return maintenancedto.CustomerMaintenanceRecordResponse{}, err
	}

	recordResponse := maintenancedto.CustomerMaintenanceRecordResponse{
		ID:                 record.ID,
		CreatedAt:          record.CreatedAt,
		Title:              record.Title,
		Details:            record.Details,
		Date:               record.CreatedAt,
		IsUserApproved:     record.IsUserApproved,
		GuaranteeViolation: violation,
	}
	return recordResponse, nil
}

func (maintenanceService *MaintenanceService) UpdateMaintenanceRequest(updateRequest maintenancedto.UpdateCustomerRequest) error {
	maintenanceRequest, err := maintenanceService.maintenanceRepository.FindRequestByID(maintenanceService.db, updateRequest.RequestID)
	if err != nil {
		return err
	}
	if maintenanceRequest == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		return notFoundError
	}

	_, err = maintenanceService.installationService.ValidatePanelOwnership(maintenanceRequest.PanelID, updateRequest.OwnerID)
	if err != nil {
		return err
	}

	if maintenanceRequest.Status != enum.MaintenanceRequestStatusPending {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.NotActive)
		return conflictErrors
	}

	if updateRequest.Subject != nil {
		maintenanceRequest.Subject = *updateRequest.Subject
	}

	if updateRequest.Description != nil {
		maintenanceRequest.Description = *updateRequest.Description
	}

	if updateRequest.UrgencyLevel != nil {
		maintenanceRequest.UrgencyLevel = enum.UrgencyLevel(*updateRequest.UrgencyLevel)
	}

	if updateRequest.IsUsingGuarantee != nil {
		if *updateRequest.IsUsingGuarantee {
			if err := maintenanceService.installationService.ValidatePanelGuarantee(maintenanceRequest.PanelID); err != nil {
				return err
			}
			maintenanceRequest.IsGuaranteeRequested = true
		} else {
			maintenanceRequest.IsGuaranteeRequested = false
		}

	}

	err = maintenanceService.maintenanceRepository.UpdateMaintenanceRequest(maintenanceService.db, maintenanceRequest)
	if err != nil {
		return err
	}
	return nil
}

func (maintenanceService *MaintenanceService) CancelMaintenanceRequest(maintenanceInfo maintenancedto.CustomerMaintenanceRequest) error {
	maintenanceRequest, err := maintenanceService.maintenanceRepository.FindRequestByID(maintenanceService.db, maintenanceInfo.RequestID)
	if err != nil {
		return err
	}
	if maintenanceRequest == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		return notFoundError
	}

	_, err = maintenanceService.installationService.ValidatePanelOwnership(maintenanceRequest.PanelID, maintenanceInfo.OwnerID)
	if err != nil {
		return err
	}

	var conflictErrors exception.ConflictErrors
	if maintenanceRequest.Status == enum.MaintenanceRequestStatusCanceled {
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.AlreadyCanceled)
		return conflictErrors
	} else if maintenanceRequest.Status != enum.MaintenanceRequestStatusPending {
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.NotActive)
		return conflictErrors
	}

	maintenanceRequest.Status = enum.MaintenanceRequestStatusCanceled

	err = maintenanceService.maintenanceRepository.UpdateMaintenanceRequest(maintenanceService.db, maintenanceRequest)
	if err != nil {
		return err
	}
	return nil
}

func (maintenanceService *MaintenanceService) ApproveMaintenanceRecord(maintenanceInfo maintenancedto.CustomerMaintenanceRequest) error {
	maintenanceRequest, err := maintenanceService.maintenanceRepository.FindRequestByID(maintenanceService.db, maintenanceInfo.RequestID)
	if err != nil {
		return err
	}
	if maintenanceRequest == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		return notFoundError
	}

	_, err = maintenanceService.installationService.ValidatePanelOwnership(maintenanceRequest.PanelID, maintenanceInfo.OwnerID)
	if err != nil {
		return err
	}

	record, err := maintenanceService.maintenanceRepository.FindRecordByRequestID(maintenanceService.db, maintenanceInfo.RequestID)
	if err != nil {
		return err
	}
	if record == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRecord}
		return notFoundError
	}

	if record.IsUserApproved {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRecord, maintenanceService.constants.Tag.AlreadyAccepted)
		return conflictErrors
	}

	record.IsUserApproved = true

	err = maintenanceService.maintenanceRepository.UpdateMaintenanceRecord(maintenanceService.db, record)
	if err != nil {
		return err
	}
	return nil
}

func (maintenanceService *MaintenanceService) GetCorporationMaintenanceRequests(listInfo maintenancedto.CorporationMaintenanceListRequest) ([]maintenancedto.CorporationMaintenanceListResponse, error) {
	err := maintenanceService.corporationService.CheckApplicantAccess(listInfo.CorporationID, listInfo.OperatorID)
	if err != nil {
		return nil, err
	}

	allowedStatus := maintenanceService.mapStatusForRole(listInfo.Status, enum.AgentTypeCorporation)

	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	maintenanceRequests, err := maintenanceService.maintenanceRepository.FindCorporationRequestsByStatus(maintenanceService.db, listInfo.CorporationID, allowedStatus, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	response := make([]maintenancedto.CorporationMaintenanceListResponse, len(maintenanceRequests))

	for i, maintenanceRequest := range maintenanceRequests {
		panelInfoRequest := installationdto.CorporationPanelRequest{
			CorporationID:  listInfo.CorporationID,
			OperatorID:     listInfo.OperatorID,
			InstallationID: maintenanceRequest.PanelID,
		}
		panel, err := maintenanceService.installationService.GetCorporationPanel(panelInfoRequest)
		if err != nil {
			return nil, err
		}

		response[i] = maintenancedto.CorporationMaintenanceListResponse{
			ID:                   maintenanceRequest.ID,
			CreatedAt:            maintenanceRequest.CreatedAt,
			Panel:                panel,
			Subject:              maintenanceRequest.Subject,
			Description:          maintenanceRequest.Description,
			UrgencyLevel:         maintenanceRequest.UrgencyLevel.String(),
			Status:               maintenanceRequest.Status.String(),
			IsGuaranteeRequested: maintenanceRequest.IsGuaranteeRequested,
		}
	}
	return response, nil
}

func (maintenanceService *MaintenanceService) GetCorporationMaintenanceRequest(maintenanceInfo maintenancedto.CorporationMaintenanceRequest) (maintenancedto.CorporationMaintenanceResponse, error) {
	err := maintenanceService.corporationService.CheckApplicantAccess(maintenanceInfo.CorporationID, maintenanceInfo.OperatorID)
	if err != nil {
		return maintenancedto.CorporationMaintenanceResponse{}, err
	}

	allowedStatuses := enum.GetAllowedMaintenanceRequestStatuses(enum.AgentTypeCorporation)
	maintenanceRequest, err := maintenanceService.maintenanceRepository.FindCorporationRequestByStatus(maintenanceService.db, maintenanceInfo.RequestID, maintenanceInfo.CorporationID, allowedStatuses)
	if err != nil {
		return maintenancedto.CorporationMaintenanceResponse{}, err
	}
	if maintenanceRequest == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		return maintenancedto.CorporationMaintenanceResponse{}, notFoundError
	}

	panelInfoRequest := installationdto.CorporationPanelRequest{
		CorporationID:  maintenanceInfo.CorporationID,
		OperatorID:     maintenanceInfo.OperatorID,
		InstallationID: maintenanceRequest.PanelID,
	}
	panel, err := maintenanceService.installationService.GetCorporationPanel(panelInfoRequest)
	if err != nil {
		return maintenancedto.CorporationMaintenanceResponse{}, err
	}

	record, err := maintenanceService.getCorporationMaintenanceRecord(maintenanceInfo.RequestID, maintenanceRequest.PanelID)
	if err != nil {
		return maintenancedto.CorporationMaintenanceResponse{}, err
	}
	response := maintenancedto.CorporationMaintenanceResponse{
		ID:                   maintenanceRequest.ID,
		CreatedAt:            maintenanceRequest.CreatedAt,
		Panel:                panel,
		Subject:              maintenanceRequest.Subject,
		Description:          maintenanceRequest.Description,
		UrgencyLevel:         maintenanceRequest.UrgencyLevel.String(),
		Status:               maintenanceRequest.Status.String(),
		IsGuaranteeRequested: maintenanceRequest.IsGuaranteeRequested,
		Record:               record,
	}
	return response, nil
}

func (maintenanceService *MaintenanceService) getCorporationMaintenanceRecord(requestID, panelID uint) (maintenancedto.CorporationMaintenanceRecordResponse, error) {
	record, err := maintenanceService.maintenanceRepository.FindRecordByRequestID(maintenanceService.db, requestID)
	if err != nil {
		return maintenancedto.CorporationMaintenanceRecordResponse{}, err
	}
	if record == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRecord}
		return maintenancedto.CorporationMaintenanceRecordResponse{}, notFoundError
	}

	operator, err := maintenanceService.userService.GetUserCredential(record.OperatorID)
	if err != nil {
		return maintenancedto.CorporationMaintenanceRecordResponse{}, err
	}

	violation, err := maintenanceService.guaranteeService.GetCorporationPanelGuaranteeViolation(panelID)
	if err != nil {
		return maintenancedto.CorporationMaintenanceRecordResponse{}, err
	}

	recordResponse := maintenancedto.CorporationMaintenanceRecordResponse{
		ID:                 record.ID,
		CreatedAt:          record.CreatedAt,
		Operator:           operator,
		Title:              record.Title,
		Details:            record.Details,
		IsUserApproved:     record.IsUserApproved,
		GuaranteeViolation: violation,
	}
	return recordResponse, nil
}

// TODO: CHECKED COULD BE BETTER add timer
func (maintenanceService *MaintenanceService) AcceptMaintenanceRequest(maintenanceInfo maintenancedto.CorporationMaintenanceRequest) error {
	err := maintenanceService.corporationService.CheckApplicantAccess(maintenanceInfo.CorporationID, maintenanceInfo.OperatorID)
	if err != nil {
		return err
	}

	allowedStatuses := enum.GetAllowedMaintenanceRequestStatuses(enum.AgentTypeCorporation)
	maintenanceRequest, err := maintenanceService.maintenanceRepository.FindCorporationRequestByStatus(maintenanceService.db, maintenanceInfo.RequestID, maintenanceInfo.CorporationID, allowedStatuses)
	if err != nil {
		return err
	}
	if maintenanceRequest == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		return notFoundError
	}

	var conflictErrors exception.ConflictErrors
	if maintenanceRequest.Status == enum.MaintenanceRequestStatusAccepted {
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.AlreadyAccepted)
		return conflictErrors
	} else if maintenanceRequest.Status == enum.MaintenanceRequestStatusRejected {
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.AlreadyRejected)
		return conflictErrors
	}

	maintenanceRequest.Status = enum.MaintenanceRequestStatusAccepted

	if err := maintenanceService.maintenanceRepository.UpdateMaintenanceRequest(maintenanceService.db, maintenanceRequest); err != nil {
		return err
	}
	return nil
}

// TODO: CHECKED COULD BE BETTER add reason
func (maintenanceService *MaintenanceService) RejectMaintenanceRequest(maintenanceInfo maintenancedto.CorporationMaintenanceRequest) error {
	err := maintenanceService.corporationService.CheckApplicantAccess(maintenanceInfo.CorporationID, maintenanceInfo.OperatorID)
	if err != nil {
		return err
	}

	allowedStatuses := enum.GetAllowedMaintenanceRequestStatuses(enum.AgentTypeCorporation)
	maintenanceRequest, err := maintenanceService.maintenanceRepository.FindCorporationRequestByStatus(maintenanceService.db, maintenanceInfo.RequestID, maintenanceInfo.CorporationID, allowedStatuses)
	if err != nil {
		return err
	}
	if maintenanceRequest == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		return notFoundError
	}

	var conflictErrors exception.ConflictErrors
	if maintenanceRequest.Status == enum.MaintenanceRequestStatusRejected {
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.AlreadyRejected)
		return conflictErrors
	} else if maintenanceRequest.Status == enum.MaintenanceRequestStatusAccepted {
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.AlreadyAccepted)
		return conflictErrors
	}

	maintenanceRequest.Status = enum.MaintenanceRequestStatusRejected

	if err := maintenanceService.maintenanceRepository.UpdateMaintenanceRequest(maintenanceService.db, maintenanceRequest); err != nil {
		return err
	}
	return nil
}

func (maintenanceService *MaintenanceService) CreateMaintenanceRecord(recordInfo maintenancedto.CreateMaintenanceRecordRequest) error {
	err := maintenanceService.corporationService.CheckApplicantAccess(recordInfo.CorporationID, recordInfo.OperatorID)
	if err != nil {
		return err
	}

	allowedStatuses := enum.GetAllowedMaintenanceRequestStatuses(enum.AgentTypeCorporation)
	maintenanceRequest, err := maintenanceService.maintenanceRepository.FindCorporationRequestByStatus(maintenanceService.db, recordInfo.RequestID, recordInfo.CorporationID, allowedStatuses)
	if err != nil {
		return err
	}
	if maintenanceRequest == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		return notFoundError
	}

	record, err := maintenanceService.maintenanceRepository.FindRecordByRequestID(maintenanceService.db, recordInfo.RequestID)
	if err != nil {
		return err
	}
	if record != nil {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRecord, maintenanceService.constants.Tag.AlreadyExist)
		return conflictErrors
	}

	var guaranteeViolationID *uint = nil
	if recordInfo.GuaranteeViolation != nil {
		recordInfo.GuaranteeViolation.PanelID = maintenanceRequest.PanelID
		request := installationdto.CreateViolatePanelGuaranteeRequest{
			CorporationID:      recordInfo.CorporationID,
			OperatorID:         recordInfo.OperatorID,
			PanelID:            maintenanceRequest.PanelID,
			GuaranteeViolation: *recordInfo.GuaranteeViolation,
		}
		temp, err := maintenanceService.installationService.ViolatePanelGuaranteeStatus(request)
		if err != nil {
			return err
		}
		guaranteeViolationID = &temp
	}

	record = &entity.MaintenanceRecord{
		OperatorID:           recordInfo.OperatorID,
		RequestID:            recordInfo.RequestID,
		IsUserApproved:       false,
		Title:                recordInfo.Title,
		Details:              recordInfo.Details,
		GuaranteeViolationID: guaranteeViolationID,
	}
	if err := maintenanceService.maintenanceRepository.CreateMaintenanceRecord(maintenanceService.db, record); err != nil {
		return err
	}

	maintenanceRequest.Status = enum.MaintenanceRequestStatusCompleted
	if err := maintenanceService.maintenanceRepository.UpdateMaintenanceRequest(maintenanceService.db, maintenanceRequest); err != nil {
		return err
	}
	return nil
}

func (maintenanceService *MaintenanceService) UpdateMaintenanceRecord(recordInfo maintenancedto.UpdateMaintenanceRecordRequest) error {
	err := maintenanceService.corporationService.CheckApplicantAccess(recordInfo.CorporationID, recordInfo.OperatorID)
	if err != nil {
		return err
	}

	allowedStatuses := enum.GetAllowedMaintenanceRequestStatuses(enum.AgentTypeCorporation)
	maintenanceRequest, err := maintenanceService.maintenanceRepository.FindCorporationRequestByStatus(maintenanceService.db, recordInfo.RequestID, recordInfo.CorporationID, allowedStatuses)
	if err != nil {
		return err
	}
	if maintenanceRequest == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		return notFoundError
	}

	maintenanceRecord, err := maintenanceService.maintenanceRepository.FindRecordByRequestID(maintenanceService.db, recordInfo.RequestID)
	if err != nil {
		return err
	}
	if maintenanceRecord == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRecord}
		return notFoundError
	}

	if maintenanceRecord.IsUserApproved {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRecord, maintenanceService.constants.Tag.NotActive)
		return conflictErrors
	}

	if recordInfo.GuaranteeViolation != nil {
		recordInfo.GuaranteeViolation.PanelID = maintenanceRequest.PanelID
		maintenanceService.guaranteeService.UpdateGuaranteeViolation(*recordInfo.GuaranteeViolation)
	}

	if recordInfo.Title != nil {
		maintenanceRecord.Title = *recordInfo.Title
	}
	if recordInfo.Details != nil {
		maintenanceRecord.Details = *recordInfo.Details
	}

	if err := maintenanceService.maintenanceRepository.UpdateMaintenanceRecord(maintenanceService.db, maintenanceRecord); err != nil {
		return err
	}
	return nil
}
