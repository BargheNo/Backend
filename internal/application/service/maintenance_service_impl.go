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
	addressService        service.AddressService
	maintenanceRepository repository.MaintenanceRepository
	db                    database.Database
}

func NewMaintenanceService(
	constants *bootstrap.Constants,
	userService service.UserService,
	installationService service.InstallationService,
	corporationService service.CorporationService,
	addressService service.AddressService,
	maintenanceRepository repository.MaintenanceRepository,
	db database.Database,
) *MaintenanceService {
	return &MaintenanceService{
		constants:             constants,
		userService:           userService,
		installationService:   installationService,
		corporationService:    corporationService,
		addressService:        addressService,
		maintenanceRepository: maintenanceRepository,
		db:                    db,
	}
}

func (maintenanceService *MaintenanceService) CreateMaintenanceRequest(requestInfo maintenancedto.NewMaintenanceRequest) {
	var conflictErrors exception.ConflictErrors
	exist := maintenanceService.userService.IsUserActive(requestInfo.OwnerID)
	if !exist {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: maintenanceService.constants.Field.MaintenanceRequest,
		}
		panic(forbiddenError)
	}
	maintenanceService.corporationService.GetCorporationByID(requestInfo.CorporationID)
	panel := maintenanceService.installationService.GetPanel(requestInfo.PanelID)

	if panel.CustomerID != requestInfo.OwnerID {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: maintenanceService.constants.Field.Panel,
		}
		panic(forbiddenError)
	}

	maintenanceRequests := maintenanceService.maintenanceRepository.FindRequestsByPanelID(maintenanceService.db, requestInfo.PanelID)
	for _, request := range maintenanceRequests {
		if request.Status == enum.MaintenanceRequestStatusPending {
			conflictErrors.Add(maintenanceService.constants.Field.MaintenanceRequest, maintenanceService.constants.Tag.Pending)
			panic(conflictErrors)
		}
	}

	maintenanceRequest := &entity.MaintenanceRequest{
		OwnerID:       requestInfo.OwnerID,
		CorporationID: requestInfo.CorporationID,
		PanelID:       requestInfo.PanelID,
		Subject:       requestInfo.Subject,
		Description:   requestInfo.Description,
		Status:        enum.MaintenanceRequestStatusPending,
		UrgencyLevel:  requestInfo.UrgencyLevel,
	}
	err := maintenanceService.maintenanceRepository.CreateMaintenanceRequest(maintenanceService.db, maintenanceRequest)
	if err != nil {
		panic(err)
	}
}

func (maintenanceService *MaintenanceService) GetCustomerMaintenanceRequests(listInfo maintenancedto.MaintenanceListRequest) []maintenancedto.MaintenanceResponse {
	maintenanceService.userService.DoesUserExist(listInfo.OwnerID)
	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	maintenanceRequests := maintenanceService.maintenanceRepository.FindMaintenanceRequestsByOwnerID(maintenanceService.db, listInfo.OwnerID, paginationModifier, sortingModifier)
	response := make([]maintenancedto.MaintenanceResponse, len(maintenanceRequests))
	for i, request := range maintenanceRequests {
		panel := maintenanceService.installationService.GetPanel(request.PanelID)
		address := maintenanceService.addressService.GetAddress(panel.ID, maintenanceService.constants.AddressOwners.Panel)
		response[i] = maintenancedto.MaintenanceResponse{
			ID:            request.ID,
			PanelID:       request.PanelID,
			CorporationID: request.CorporationID,
			OwnerID:       request.OwnerID,
			Subject:       request.Subject,
			Description:   request.Description,
			UrgencyLevel:  request.UrgencyLevel.String(),
			Status:        request.Status.String(),
			CreatedAt:     request.CreatedAt,
			Panel: installationdto.CustomerPanelResponse{
				ID:                   panel.ID,
				PanelName:            panel.Name,
				Power:                panel.Power,
				Area:                 panel.Area,
				BuildingType:         panel.BuildingType,
				Tilt:                 panel.Tilt,
				Azimuth:              panel.Azimuth,
				TotalNumberOfModules: panel.TotalNumberOfModules,
				Address:              address,
				CorporationName:      panel.Corporation.Name,
			},
		}
	}
	return response
}

func (maintenanceService *MaintenanceService) GetCorporationMaintenanceRequests(listInfo maintenancedto.CorporationMaintenanceListRequest) []maintenancedto.CorporationMaintenanceResponse {
	maintenanceService.corporationService.CheckApplicantAccess(listInfo.CorporationID, listInfo.OperatorID)
	maintenanceService.corporationService.GetCorporationByID(listInfo.CorporationID)
	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	maintenanceRequests := maintenanceService.maintenanceRepository.FindMaintenanceRequestsByCorporationID(maintenanceService.db, listInfo.CorporationID, paginationModifier, sortingModifier)
	response := make([]maintenancedto.CorporationMaintenanceResponse, len(maintenanceRequests))
	for i, request := range maintenanceRequests {
		panel := maintenanceService.installationService.GetPanel(request.PanelID)
		address := maintenanceService.addressService.GetAddress(panel.ID, maintenanceService.constants.AddressOwners.Panel)
		response[i] = maintenancedto.CorporationMaintenanceResponse{
			ID:           request.ID,
			PanelID:      request.PanelID,
			Subject:      request.Subject,
			Description:  request.Description,
			UrgencyLevel: request.UrgencyLevel.String(),
			Status:       request.Status.String(),
			CreatedAt:    request.CreatedAt,
			OwnerPhone:   panel.Customer.Phone,
			Panel: installationdto.CorporationPanelResponse{
				ID:                   panel.ID,
				PanelName:            panel.Name,
				Power:                panel.Power,
				Area:                 panel.Area,
				BuildingType:         panel.BuildingType,
				Tilt:                 panel.Tilt,
				Azimuth:              panel.Azimuth,
				TotalNumberOfModules: panel.TotalNumberOfModules,
				Address:              address,
				OperatorName:         panel.Operator.FirstName + " " + panel.Operator.LastName,
			},
		}
	}
	return response
}

func (maintenanceService *MaintenanceService) HandleRequest(handleRequestInfo maintenancedto.HandleRequest) {
	maintenanceService.corporationService.CheckApplicantAccess(handleRequestInfo.CorporationID, handleRequestInfo.OperatorID)
	request := maintenanceService.maintenanceRepository.FindMaintenanceRequestByID(maintenanceService.db, handleRequestInfo.RequestID)
	if request == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		panic(notFoundError)
	}
	if request.Status != enum.MaintenanceRequestStatusPending {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: maintenanceService.constants.Field.MaintenanceRequest,
		}
		panic(forbiddenError)
	}
	if handleRequestInfo.Accept {
		request.Status = enum.MaintenanceRequestStatusAccepted
	} else {
		request.Status = enum.MaintenanceRequestStatusRejected
	}
	err := maintenanceService.maintenanceRepository.UpdateMaintenanceRequest(maintenanceService.db, request)
	if err != nil {
		panic(err)
	}
}

func (maintenanceService *MaintenanceService) AddMaintenanceRecord(requestInfo maintenancedto.AddMaintenanceRecordRequest) {
	maintenanceService.corporationService.CheckApplicantAccess(requestInfo.CorporationID, requestInfo.OperatorID)
	request := maintenanceService.maintenanceRepository.FindMaintenanceRequestByID(maintenanceService.db, requestInfo.RequestID)
	if request == nil {
		notFoundError := exception.NotFoundError{Item: maintenanceService.constants.Field.MaintenanceRequest}
		panic(notFoundError)
	}
	if request.Status != enum.MaintenanceRequestStatusAccepted {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: maintenanceService.constants.Field.MaintenanceRequest,
		}
		panic(forbiddenError)
	}
	record := &entity.MaintenanceRecord{
		PanelID:       request.PanelID,
		CustomerID:    request.OwnerID,
		CorporationID: request.CorporationID,
		OperatorID:    requestInfo.OperatorID,
		Title:         requestInfo.Title,
		Details:       requestInfo.Details,
		Date:          requestInfo.Date,
	}
	request.Status = enum.MaintenanceRequestStatusCompleted
	err := maintenanceService.maintenanceRepository.CreateMaintenanceRecord(maintenanceService.db, record)
	if err != nil {
		panic(err)
	}
	err = maintenanceService.maintenanceRepository.UpdateMaintenanceRequest(maintenanceService.db, request)
	if err != nil {
		panic(err)
	}
}

func (maintenanceService *MaintenanceService) GetCorporationMaintenanceRecords(requestInfo maintenancedto.CorporationMaintenanceListRequest) []maintenancedto.MaintenanceRecordResponse {
	maintenanceService.corporationService.CheckApplicantAccess(requestInfo.CorporationID, requestInfo.OperatorID)
	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	maintenanceRecords := maintenanceService.maintenanceRepository.FindMaintenanceRecordsByCorporationID(maintenanceService.db, requestInfo.CorporationID, paginationModifier, sortingModifier)
	response := make([]maintenancedto.MaintenanceRecordResponse, len(maintenanceRecords))
	for i, record := range maintenanceRecords {
		panel := maintenanceService.installationService.GetPanel(record.PanelID)
		address := maintenanceService.addressService.GetAddress(panel.ID, maintenanceService.constants.AddressOwners.Panel)
		response[i] = maintenancedto.MaintenanceRecordResponse{
			ID:        record.ID,
			RequestID: record.ID,
			Panel: installationdto.CorporationPanelResponse{
				ID:                   panel.ID,
				PanelName:            panel.Name,
				Power:                panel.Power,
				Area:                 panel.Area,
				BuildingType:         panel.BuildingType,
				Tilt:                 panel.Tilt,
				Azimuth:              panel.Azimuth,
				TotalNumberOfModules: panel.TotalNumberOfModules,
				Address:              address,
				OperatorName:         panel.Operator.FirstName + " " + panel.Operator.LastName,
			},
			OperatorID:    requestInfo.OperatorID,
			CorporationID: requestInfo.CorporationID,
			Title:         record.Title,
			Details:       record.Details,
			Date:          record.Date,
		}
	}
	return response
}

func (maintenanceService *MaintenanceService) GetCorporationMaintenanceRecordsByPanel(requestInfo maintenancedto.CorporationMaintenanceRecordByPanelRequest) []maintenancedto.MaintenanceRecordResponse {
	maintenanceService.corporationService.CheckApplicantAccess(requestInfo.CorporationID, requestInfo.OperatorID)
	maintenanceService.installationService.GetPanel(requestInfo.PanelID)
	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	maintenanceRecords := maintenanceService.maintenanceRepository.FindMaintenanceRecordsByPanelAndCorporationID(maintenanceService.db, requestInfo.PanelID, requestInfo.CorporationID, paginationModifier, sortingModifier)
	response := make([]maintenancedto.MaintenanceRecordResponse, len(maintenanceRecords))
	for i, record := range maintenanceRecords {
		panel := maintenanceService.installationService.GetPanel(record.PanelID)
		address := maintenanceService.addressService.GetAddress(panel.ID, maintenanceService.constants.AddressOwners.Panel)
		response[i] = maintenancedto.MaintenanceRecordResponse{
			ID:        record.ID,
			RequestID: record.ID,
			Panel: installationdto.CorporationPanelResponse{
				ID:                   panel.ID,
				PanelName:            panel.Name,
				Power:                panel.Power,
				Area:                 panel.Area,
				BuildingType:         panel.BuildingType,
				Tilt:                 panel.Tilt,
				Azimuth:              panel.Azimuth,
				TotalNumberOfModules: panel.TotalNumberOfModules,
				Address:              address,
				OperatorName:         panel.Operator.FirstName + " " + panel.Operator.LastName,
			},
			OperatorID:    requestInfo.OperatorID,
			CorporationID: requestInfo.CorporationID,
			Title:         record.Title,
			Details:       record.Details,
			Date:          record.Date,
		}
	}
	return response
}

func (maintenanceService *MaintenanceService) GetCustomerMaintenanceRecords(requestInfo maintenancedto.MaintenanceListRequest) []maintenancedto.CustomerMaintenanceRecordResponse {
	maintenanceService.userService.DoesUserExist(requestInfo.OwnerID)
	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	maintenanceRecords := maintenanceService.maintenanceRepository.FindMaintenanceRecordsByCustomerID(maintenanceService.db, requestInfo.OwnerID, paginationModifier, sortingModifier)
	response := make([]maintenancedto.CustomerMaintenanceRecordResponse, len(maintenanceRecords))
	for i, record := range maintenanceRecords {
		panel := maintenanceService.installationService.GetPanel(record.PanelID)
		address := maintenanceService.addressService.GetAddress(panel.ID, maintenanceService.constants.AddressOwners.Panel)
		response[i] = maintenancedto.CustomerMaintenanceRecordResponse{
			ID: record.ID,
			Panel: installationdto.CustomerPanelResponse{
				ID:                   panel.ID,
				PanelName:            panel.Name,
				Power:                panel.Power,
				Area:                 panel.Area,
				BuildingType:         panel.BuildingType,
				Tilt:                 panel.Tilt,
				Azimuth:              panel.Azimuth,
				TotalNumberOfModules: panel.TotalNumberOfModules,
				Address:              address,
				CorporationName:      panel.Corporation.Name,
			},
			OperatorID:    record.OperatorID,
			OperatorPhone: record.Operator.Phone,
			Title:         record.Title,
			Details:       record.Details,
			Date:          record.Date,
		}
	}
	return response
}

func (maintenanceService *MaintenanceService) GetCustomerMaintenanceRecordsByPanel(requestInfo maintenancedto.CustomerMaintenanceRecordByPanelRequest) []maintenancedto.CustomerMaintenanceRecordResponse {
	maintenanceService.userService.DoesUserExist(requestInfo.OwnerID)
	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	maintenanceRecords := maintenanceService.maintenanceRepository.FindCustomerMaintenanceRecordsByPanelID(maintenanceService.db, requestInfo.PanelID, requestInfo.OwnerID, paginationModifier, sortingModifier)
	response := make([]maintenancedto.CustomerMaintenanceRecordResponse, len(maintenanceRecords))
	for i, record := range maintenanceRecords {
		panel := maintenanceService.installationService.GetPanel(record.PanelID)
		address := maintenanceService.addressService.GetAddress(panel.ID, maintenanceService.constants.AddressOwners.Panel)
		response[i] = maintenancedto.CustomerMaintenanceRecordResponse{
			ID: record.ID,
			Panel: installationdto.CustomerPanelResponse{
				ID:                   panel.ID,
				PanelName:            panel.Name,
				Power:                panel.Power,
				Area:                 panel.Area,
				BuildingType:         panel.BuildingType,
				Tilt:                 panel.Tilt,
				Azimuth:              panel.Azimuth,
				TotalNumberOfModules: panel.TotalNumberOfModules,
				Address:              address,
				CorporationName:      panel.Corporation.Name,
			},
			OperatorID:    record.OperatorID,
			OperatorPhone: record.Operator.Phone,
			Title:         record.Title,
			Details:       record.Details,
			Date:          record.Date,
		}
	}
	return response
}
