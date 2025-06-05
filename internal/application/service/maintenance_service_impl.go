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
	maintenanceRecord, exist := maintenanceService.maintenanceRepository.FindRecordByID(maintenanceService.db, recordID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRecord}
		return notFoundError
	}

	maintenanceRequest, exist := maintenanceService.maintenanceRepository.FindRequestByID(maintenanceService.db, maintenanceRecord.RequestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		return notFoundError
	}

	if err := maintenanceService.installationService.ValidatePanelOwnership(maintenanceRequest.PanelID, userID); err != nil {
		return err
	}

	return nil
}

func (maintenanceService *MaintenanceService) GetRequestByAdmin(recordID uint) (maintenancedto.AdminMaintenanceRequestResponse, error) {
	var response maintenancedto.AdminMaintenanceRequestResponse
	maintenanceRecord, exist := maintenanceService.maintenanceRepository.FindRecordByID(maintenanceService.db, recordID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRecord}
		return response, notFoundError
	}

	maintenanceRequest, exist := maintenanceService.maintenanceRepository.FindRequestByID(maintenanceService.db, maintenanceRecord.RequestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		return response, notFoundError
	}

	panel := maintenanceService.installationService.GetPanelByAdmin(maintenanceRequest.PanelID)

	corporation := maintenanceService.corporationService.GetCorporationCredentials(maintenanceRequest.CorporationID)

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

// CHECKED
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

// CHECKED
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

// CHECKED
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

// CHECKED
func (maintenanceService *MaintenanceService) CreateMaintenanceRequest(request maintenancedto.CreateMaintenanceRequest) {
	if exist := maintenanceService.userService.IsUserActive(request.OwnerID); !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: maintenanceService.constants.Field.MaintenanceRequest,
		}
		panic(forbiddenError)
	}

	maintenanceService.corporationService.DoesCorporationExist(request.CorporationID)

	if err := maintenanceService.installationService.ValidatePanelOwnership(request.PanelID, request.OwnerID); err != nil {
		panic(err)
	}

	allowedStatus := []enum.MaintenanceRequestStatus{enum.MaintenanceRequestStatusPending}
	currentActiveRequests := maintenanceService.maintenanceRepository.FindRequestsByPanelIDAndStatus(maintenanceService.db, request.PanelID, allowedStatus)
	if len(currentActiveRequests) != 0 {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.Pending)
		panic(conflictErrors)
	}

	if request.IsUsingGuarantee {
		if err := maintenanceService.installationService.ValidatePanelGuarantee(request.PanelID); err != nil {
			panic(err)
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
	err := maintenanceService.maintenanceRepository.CreateMaintenanceRequest(maintenanceService.db, maintenanceRequest)
	if err != nil {
		panic(err)
	}
}

// CHECKED
func (maintenanceService *MaintenanceService) GetCustomerMaintenanceRequests(listInfo maintenancedto.CustomerMaintenanceListRequest) []maintenancedto.CustomerMaintenanceRequestResponse {
	allowedStatus := maintenanceService.mapStatusForRole(listInfo.Status, enum.AgentTypeCustomer)

	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	maintenanceRequests := maintenanceService.maintenanceRepository.FindRequestsByCustomerID(maintenanceService.db, listInfo.OwnerID, allowedStatus, paginationModifier, sortingModifier)
	response := make([]maintenancedto.CustomerMaintenanceRequestResponse, len(maintenanceRequests))

	for i, maintenanceRequest := range maintenanceRequests {
		panelInfoRequest := installationdto.GetOwnerRequest{
			OwnerID:        listInfo.OwnerID,
			InstallationID: maintenanceRequest.PanelID,
		}
		panel := maintenanceService.installationService.GetCustomerPanel(panelInfoRequest)

		corporation := maintenanceService.corporationService.GetCorporationCredentials(maintenanceRequest.CorporationID)

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
	return response
}

// CHECKED
func (maintenanceService *MaintenanceService) GetCustomerPanelMaintenanceRequests(listInfo maintenancedto.CustomerPanelMaintenanceListRequest) []maintenancedto.CustomerMaintenanceRequestResponse {
	if err := maintenanceService.installationService.ValidatePanelOwnership(listInfo.PanelID, listInfo.OwnerID); err != nil {
		panic(err)
	}

	allowedStatus := maintenanceService.mapStatusForRole(listInfo.Status, enum.AgentTypeCustomer)

	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	maintenanceRequests := maintenanceService.maintenanceRepository.FindRequestsByPanelIDAndStatus(maintenanceService.db, listInfo.PanelID, allowedStatus, paginationModifier, sortingModifier)
	response := make([]maintenancedto.CustomerMaintenanceRequestResponse, len(maintenanceRequests))

	for i, maintenanceRequest := range maintenanceRequests {
		panelInfoRequest := installationdto.GetOwnerRequest{
			OwnerID:        listInfo.OwnerID,
			InstallationID: maintenanceRequest.PanelID,
		}
		panel := maintenanceService.installationService.GetCustomerPanel(panelInfoRequest)

		corporation := maintenanceService.corporationService.GetCorporationCredentials(maintenanceRequest.CorporationID)

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
	return response
}

// CHECKED
func (maintenanceService *MaintenanceService) GetCustomerMaintenanceRequest(maintenanceInfo maintenancedto.CustomerMaintenanceRequest) maintenancedto.CustomerMaintenanceRequestResponse {
	maintenanceRequest, exist := maintenanceService.maintenanceRepository.FindRequestByID(maintenanceService.db, maintenanceInfo.RequestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		panic(notFoundError)
	}

	panelInfoRequest := installationdto.GetOwnerRequest{
		OwnerID:        maintenanceInfo.OwnerID,
		InstallationID: maintenanceRequest.PanelID,
	}
	panel := maintenanceService.installationService.GetCustomerPanel(panelInfoRequest)

	corporation := maintenanceService.corporationService.GetCorporationCredentials(maintenanceRequest.CorporationID)

	record, _ := maintenanceService.getCustomerMaintenanceRecord(maintenanceInfo.RequestID, maintenanceRequest.PanelID)

	return maintenancedto.CustomerMaintenanceRequestResponse{
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
}

func (maintenanceService *MaintenanceService) getCustomerMaintenanceRecord(requestID, panelID uint) (maintenancedto.CustomerMaintenanceRecordResponse, error) {
	record, exist := maintenanceService.maintenanceRepository.FindRecordByRequestID(maintenanceService.db, requestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRecord}
		return maintenancedto.CustomerMaintenanceRecordResponse{}, notFoundError
	}

	violation, _ := maintenanceService.guaranteeService.GetCustomerPanelGuaranteeViolation(panelID)

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

// CHECKED
func (maintenanceService *MaintenanceService) UpdateMaintenanceRequest(updateRequest maintenancedto.UpdateCustomerRequest) {
	maintenanceRequest, exist := maintenanceService.maintenanceRepository.FindRequestByID(maintenanceService.db, updateRequest.RequestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		panic(notFoundError)
	}

	if err := maintenanceService.installationService.ValidatePanelOwnership(maintenanceRequest.PanelID, updateRequest.OwnerID); err != nil {
		panic(err)
	}

	if maintenanceRequest.Status != enum.MaintenanceRequestStatusPending {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.NotActive)
		panic(conflictErrors)

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
				panic(err)
			}
			maintenanceRequest.IsGuaranteeRequested = true
		} else {
			maintenanceRequest.IsGuaranteeRequested = false
		}

	}

	err := maintenanceService.maintenanceRepository.UpdateMaintenanceRequest(maintenanceService.db, maintenanceRequest)
	if err != nil {
		panic(err)
	}
}

// CHECKED
func (maintenanceService *MaintenanceService) CancelMaintenanceRequest(maintenanceInfo maintenancedto.CustomerMaintenanceRequest) {
	maintenanceRequest, exist := maintenanceService.maintenanceRepository.FindRequestByID(maintenanceService.db, maintenanceInfo.RequestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		panic(notFoundError)
	}

	if err := maintenanceService.installationService.ValidatePanelOwnership(maintenanceRequest.PanelID, maintenanceInfo.OwnerID); err != nil {
		panic(err)
	}

	var conflictErrors exception.ConflictErrors
	if maintenanceRequest.Status == enum.MaintenanceRequestStatusCanceled {
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.AlreadyCanceled)
		panic(conflictErrors)
	} else if maintenanceRequest.Status != enum.MaintenanceRequestStatusPending {
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.NotActive)
		panic(conflictErrors)
	}

	maintenanceRequest.Status = enum.MaintenanceRequestStatusCanceled

	err := maintenanceService.maintenanceRepository.UpdateMaintenanceRequest(maintenanceService.db, maintenanceRequest)
	if err != nil {
		panic(err)
	}
}

// CHECKED
func (maintenanceService *MaintenanceService) ApproveMaintenanceRecord(maintenanceInfo maintenancedto.CustomerMaintenanceRequest) {
	maintenanceRequest, exist := maintenanceService.maintenanceRepository.FindRequestByID(maintenanceService.db, maintenanceInfo.RequestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		panic(notFoundError)
	}

	if err := maintenanceService.installationService.ValidatePanelOwnership(maintenanceRequest.PanelID, maintenanceInfo.OwnerID); err != nil {
		panic(err)
	}

	record, exist := maintenanceService.maintenanceRepository.FindRecordByRequestID(maintenanceService.db, maintenanceInfo.RequestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRecord}
		panic(notFoundError)
	}

	if record.IsUserApproved {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRecord, maintenanceService.constants.Tag.AlreadyAccepted)
		panic(conflictErrors)
	}

	record.IsUserApproved = true

	err := maintenanceService.maintenanceRepository.UpdateMaintenanceRecord(maintenanceService.db, record)
	if err != nil {
		panic(err)
	}
}

// CHECKED
func (maintenanceService *MaintenanceService) GetCorporationMaintenanceRequests(listInfo maintenancedto.CorporationMaintenanceListRequest) []maintenancedto.CorporationMaintenanceListResponse {
	maintenanceService.corporationService.CheckApplicantAccess(listInfo.CorporationID, listInfo.OperatorID)

	allowedStatus := maintenanceService.mapStatusForRole(listInfo.Status, enum.AgentTypeCorporation)

	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	maintenanceRequests := maintenanceService.maintenanceRepository.FindCorporationRequestsByStatus(maintenanceService.db, listInfo.CorporationID, allowedStatus, paginationModifier, sortingModifier)
	response := make([]maintenancedto.CorporationMaintenanceListResponse, len(maintenanceRequests))

	for i, maintenanceRequest := range maintenanceRequests {
		panelInfoRequest := installationdto.CorporationPanelRequest{
			CorporationID:  listInfo.CorporationID,
			OperatorID:     listInfo.OperatorID,
			InstallationID: maintenanceRequest.PanelID,
		}
		panel := maintenanceService.installationService.GetCorporationPanel(panelInfoRequest)

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
	return response
}

// CHECKED
func (maintenanceService *MaintenanceService) GetCorporationMaintenanceRequest(maintenanceInfo maintenancedto.CorporationMaintenanceRequest) maintenancedto.CorporationMaintenanceResponse {
	maintenanceService.corporationService.CheckApplicantAccess(maintenanceInfo.CorporationID, maintenanceInfo.OperatorID)

	allowedStatuses := enum.GetAllowedMaintenanceRequestStatuses(enum.AgentTypeCorporation)
	maintenanceRequest, exist := maintenanceService.maintenanceRepository.FindCorporationRequestByStatus(maintenanceService.db, maintenanceInfo.RequestID, maintenanceInfo.CorporationID, allowedStatuses)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		panic(notFoundError)
	}

	panelInfoRequest := installationdto.CorporationPanelRequest{
		CorporationID:  maintenanceInfo.CorporationID,
		OperatorID:     maintenanceInfo.OperatorID,
		InstallationID: maintenanceRequest.PanelID,
	}
	panel := maintenanceService.installationService.GetCorporationPanel(panelInfoRequest)

	record, _ := maintenanceService.getCorporationMaintenanceRecord(maintenanceInfo.RequestID, maintenanceRequest.PanelID)
	return maintenancedto.CorporationMaintenanceResponse{
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
}

// CHECKED
func (maintenanceService *MaintenanceService) getCorporationMaintenanceRecord(requestID, panelID uint) (maintenancedto.CorporationMaintenanceRecordResponse, error) {
	record, exist := maintenanceService.maintenanceRepository.FindRecordByRequestID(maintenanceService.db, requestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRecord}
		return maintenancedto.CorporationMaintenanceRecordResponse{}, notFoundError
	}

	operator := maintenanceService.userService.GetUserCredential(record.OperatorID)

	violation, _ := maintenanceService.guaranteeService.GetCorporationPanelGuaranteeViolation(panelID)

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
func (maintenanceService *MaintenanceService) AcceptMaintenanceRequest(maintenanceInfo maintenancedto.CorporationMaintenanceRequest) {
	maintenanceService.corporationService.CheckApplicantAccess(maintenanceInfo.CorporationID, maintenanceInfo.OperatorID)

	allowedStatuses := enum.GetAllowedMaintenanceRequestStatuses(enum.AgentTypeCorporation)
	maintenanceRequest, exist := maintenanceService.maintenanceRepository.FindCorporationRequestByStatus(maintenanceService.db, maintenanceInfo.RequestID, maintenanceInfo.CorporationID, allowedStatuses)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		panic(notFoundError)
	}

	var conflictErrors exception.ConflictErrors
	if maintenanceRequest.Status == enum.MaintenanceRequestStatusAccepted {
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.AlreadyAccepted)
		panic(conflictErrors)
	} else if maintenanceRequest.Status == enum.MaintenanceRequestStatusRejected {
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.AlreadyRejected)
		panic(conflictErrors)
	}

	maintenanceRequest.Status = enum.MaintenanceRequestStatusAccepted

	if err := maintenanceService.maintenanceRepository.UpdateMaintenanceRequest(maintenanceService.db, maintenanceRequest); err != nil {
		panic(err)
	}
}

// TODO: CHECKED COULD BE BETTER add reason
func (maintenanceService *MaintenanceService) RejectMaintenanceRequest(maintenanceInfo maintenancedto.CorporationMaintenanceRequest) {
	maintenanceService.corporationService.CheckApplicantAccess(maintenanceInfo.CorporationID, maintenanceInfo.OperatorID)

	allowedStatuses := enum.GetAllowedMaintenanceRequestStatuses(enum.AgentTypeCorporation)
	maintenanceRequest, exist := maintenanceService.maintenanceRepository.FindCorporationRequestByStatus(maintenanceService.db, maintenanceInfo.RequestID, maintenanceInfo.CorporationID, allowedStatuses)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		panic(notFoundError)
	}

	var conflictErrors exception.ConflictErrors
	if maintenanceRequest.Status == enum.MaintenanceRequestStatusRejected {
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.AlreadyRejected)
		panic(conflictErrors)
	} else if maintenanceRequest.Status == enum.MaintenanceRequestStatusAccepted {
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.AlreadyAccepted)
		panic(conflictErrors)
	}

	maintenanceRequest.Status = enum.MaintenanceRequestStatusRejected

	if err := maintenanceService.maintenanceRepository.UpdateMaintenanceRequest(maintenanceService.db, maintenanceRequest); err != nil {
		panic(err)
	}
}

// CHECKED
func (maintenanceService *MaintenanceService) CreateMaintenanceRecord(recordInfo maintenancedto.CreateMaintenanceRecordRequest) {
	maintenanceService.corporationService.CheckApplicantAccess(recordInfo.CorporationID, recordInfo.OperatorID)

	allowedStatuses := enum.GetAllowedMaintenanceRequestStatuses(enum.AgentTypeCorporation)
	maintenanceRequest, exist := maintenanceService.maintenanceRepository.FindCorporationRequestByStatus(maintenanceService.db, recordInfo.RequestID, recordInfo.CorporationID, allowedStatuses)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		panic(notFoundError)
	}

	if _, exist = maintenanceService.maintenanceRepository.FindRecordByRequestID(maintenanceService.db, recordInfo.RequestID); exist {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRecord, maintenanceService.constants.Tag.AlreadyExist)
		panic(conflictErrors)
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
		temp := maintenanceService.installationService.ViolatePanelGuaranteeStatus(request)
		guaranteeViolationID = &temp
	}

	record := &entity.MaintenanceRecord{
		OperatorID:           recordInfo.OperatorID,
		RequestID:            recordInfo.RequestID,
		IsUserApproved:       false,
		Title:                recordInfo.Title,
		Details:              recordInfo.Details,
		GuaranteeViolationID: guaranteeViolationID,
	}
	if err := maintenanceService.maintenanceRepository.CreateMaintenanceRecord(maintenanceService.db, record); err != nil {
		panic(err)
	}

	maintenanceRequest.Status = enum.MaintenanceRequestStatusCompleted
	if err := maintenanceService.maintenanceRepository.UpdateMaintenanceRequest(maintenanceService.db, maintenanceRequest); err != nil {
		panic(err)
	}
}

// CHECKED
func (maintenanceService *MaintenanceService) UpdateMaintenanceRecord(recordInfo maintenancedto.UpdateMaintenanceRecordRequest) {
	maintenanceService.corporationService.CheckApplicantAccess(recordInfo.CorporationID, recordInfo.OperatorID)

	allowedStatuses := enum.GetAllowedMaintenanceRequestStatuses(enum.AgentTypeCorporation)
	maintenanceRequest, exist := maintenanceService.maintenanceRepository.FindCorporationRequestByStatus(maintenanceService.db, recordInfo.RequestID, recordInfo.CorporationID, allowedStatuses)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		panic(notFoundError)
	}

	maintenanceRecord, exist := maintenanceService.maintenanceRepository.FindRecordByRequestID(maintenanceService.db, recordInfo.RequestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRecord}
		panic(notFoundError)
	}

	if maintenanceRecord.IsUserApproved {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRecord, maintenanceService.constants.Tag.NotActive)
		panic(conflictErrors)
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
		panic(err)
	}
}
