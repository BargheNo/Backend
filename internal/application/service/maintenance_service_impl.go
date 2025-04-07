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
	maintenanceService.userService.GetUserCredential(requestInfo.OwnerID)
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
	maintenanceService.userService.GetUserCredential(listInfo.OwnerID)
	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	maintenanceRequests := maintenanceService.maintenanceRepository.FindMaintenanceRequestsByOwnerID(maintenanceService.db, listInfo.OwnerID, paginationModifier, sortingModifier)
	var maintenanceResponses []maintenancedto.MaintenanceResponse
	for _, request := range maintenanceRequests {
		panel := maintenanceService.installationService.GetPanel(request.PanelID)
		address := maintenanceService.addressService.GetAddress(panel.ID, maintenanceService.constants.AddressOwners.Panel)
		maintenanceResponse := maintenancedto.MaintenanceResponse{
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
		maintenanceResponses = append(maintenanceResponses, maintenanceResponse)
	}
	return maintenanceResponses
}

func (maintenanceService *MaintenanceService) GetCorporationMaintenanceRequests(listInfo maintenancedto.CorporationMaintenanceListRequest) []maintenancedto.CorporationMaintenanceResponse {
	maintenanceService.corporationService.CheckApplicantAccess(listInfo.OperatorID, listInfo.CorporationID)
	maintenanceService.corporationService.GetCorporationByID(listInfo.CorporationID)
	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	maintenanceRequests := maintenanceService.maintenanceRepository.FindMaintenanceRequestsByCorporationID(maintenanceService.db, listInfo.CorporationID, paginationModifier, sortingModifier)
	var maintenanceResponses []maintenancedto.CorporationMaintenanceResponse
	for _, request := range maintenanceRequests {
		panel := maintenanceService.installationService.GetPanel(request.PanelID)
		address := maintenanceService.addressService.GetAddress(panel.ID, maintenanceService.constants.AddressOwners.Panel)
		maintenanceResponse := maintenancedto.CorporationMaintenanceResponse{
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
		maintenanceResponses = append(maintenanceResponses, maintenanceResponse)
	}
	return maintenanceResponses
}
