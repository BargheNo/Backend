package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
)

type InstallationService struct {
	constants              *bootstrap.Constants
	addressService         service.AddressService
	userService            service.UserService
	corporationService     service.CorporationService
	installationRepository repository.InstallationRepository
	db                     database.Database
}

func NewInstallationService(
	constants *bootstrap.Constants,
	addressService service.AddressService,
	userService service.UserService,
	corporationService service.CorporationService,
	installationRepository repository.InstallationRepository,
	db database.Database,
) *InstallationService {
	return &InstallationService{
		constants:              constants,
		addressService:         addressService,
		userService:            userService,
		corporationService:     corporationService,
		installationRepository: installationRepository,
		db:                     db,
	}
}

func (installationService *InstallationService) GetInstallationRequestModel(requestID uint) *entity.InstallationRequest {
	request, exist := installationService.installationRepository.FindRequestByID(installationService.db, requestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.InstallationRequest}
		panic(notFoundError)
	}
	return request
}

func (installationService *InstallationService) CreateInstallationRequest(requestInfo installationdto.NewInstallationRequest) {
	// get user by id from user service and check complete tag and if not completed -> 403 forbidden
	// compare installed panels names to new request name
	allowedStatus := []enum.InstallationRequestStatus{enum.InstallationRequestStatusActive}
	_, exist := installationService.installationRepository.FindOwnerRequestByName(installationService.db, requestInfo.OwnerID, allowedStatus, requestInfo.Name)
	if exist {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(installationService.constants.Field.Name, installationService.constants.Tag.AlreadyRegistered)
		panic(conflictErrors)
	}
	inProgressReqs := installationService.installationRepository.FindOwnerRequests(installationService.db, requestInfo.OwnerID, allowedStatus)
	if len(inProgressReqs) >= 5 {
		rateLimitError := exception.NewConcurrentInstallLimitError("", 5, nil)
		panic(rateLimitError)
	}

	request := &entity.InstallationRequest{
		Name:         requestInfo.Name,
		Status:       enum.InstallationRequestStatusActive,
		Area:         requestInfo.Area,
		PowerRequest: requestInfo.Power,
		MaxCost:      requestInfo.MaxCost,
		BuildingType: requestInfo.BuildingType,
		OwnerID:      requestInfo.OwnerID,
		Address: entity.Address{
			ProvinceID:    requestInfo.Address.ProvinceID,
			CityID:        requestInfo.Address.CityID,
			StreetAddress: requestInfo.Address.StreetAddress,
			PostalCode:    requestInfo.Address.PostalCode,
			HouseNumber:   requestInfo.Address.HouseNumber,
			Unit:          requestInfo.Address.Unit,
		},
	}
	err := installationService.installationRepository.CreateRequest(installationService.db, request)
	if err != nil {
		panic(err)
	}
}

func (installationService *InstallationService) GetOwnerInstallationRequests(listInfo installationdto.InstallationListRequest) []installationdto.OwnerRequestsResponse {
	allowedStatus := []enum.InstallationRequestStatus{enum.InstallationRequestStatusActive, enum.InstallationRequestStatusCancelled, enum.InstallationRequestStatusExpired}
	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	requests := installationService.installationRepository.FindOwnerRequests(
		installationService.db, listInfo.OwnerID, allowedStatus, paginationModifier, sortingModifier)
	response := make([]installationdto.OwnerRequestsResponse, len(requests))

	for i, request := range requests {
		address := installationService.addressService.GetAddress(request.OwnerID, installationService.constants.AddressOwners.InstallationRequest)
		response[i] = installationdto.OwnerRequestsResponse{
			ID:           request.ID,
			Name:         request.Name,
			CreatedTime:  request.CreatedAt,
			Status:       request.Status.String(),
			PowerRequest: request.PowerRequest,
			MaxCost:      request.MaxCost,
			BuildingType: request.BuildingType,
			Address:      address,
		}
	}
	return response
}

func (installationService *InstallationService) GetInstallationRequest(requestID uint) installationdto.RequestDetailsResponse {
	request := installationService.GetInstallationRequestModel(requestID)
	address := installationService.addressService.GetAddress(request.OwnerID, installationService.constants.AddressOwners.InstallationRequest)
	customer := installationService.userService.GetUserCredential(request.OwnerID)
	return installationdto.RequestDetailsResponse{
		ID:           request.ID,
		Name:         request.Name,
		CreatedTime:  request.CreatedAt,
		Status:       request.Status.String(),
		PowerRequest: request.PowerRequest,
		MaxCost:      request.MaxCost,
		BuildingType: request.BuildingType,
		Address:      address,
		Customer:     customer,
	}
}

func (installationService *InstallationService) GetOwnerInstallationRequest(requestInfo installationdto.GetOwnerRequest) installationdto.OwnerRequestsResponse {
	installationRequest := installationService.GetInstallationRequestModel(requestInfo.RequestID)
	if installationRequest.OwnerID != requestInfo.OwnerID {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: installationService.constants.Field.InstallationRequest,
		}
		panic(forbiddenError)
	}
	address := installationService.addressService.GetAddress(requestInfo.OwnerID, installationService.constants.AddressOwners.InstallationRequest)

	return installationdto.OwnerRequestsResponse{
		ID:           installationRequest.ID,
		Name:         installationRequest.Name,
		CreatedTime:  installationRequest.CreatedAt,
		Status:       installationRequest.Status.String(),
		PowerRequest: installationRequest.PowerRequest,
		MaxCost:      installationRequest.MaxCost,
		BuildingType: installationRequest.BuildingType,
		Address:      address,
	}
}

func (installationService *InstallationService) GetInstallationRequests(listInfo installationdto.InstallationListRequest) []installationdto.RequestDetailsResponse {
	allowedStatus := []enum.InstallationRequestStatus{enum.InstallationRequestStatusActive}
	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	requests := installationService.installationRepository.FindRequestByStatus(installationService.db, allowedStatus, paginationModifier, sortingModifier)
	response := make([]installationdto.RequestDetailsResponse, len(requests))
	for i, request := range requests {
		address := installationService.addressService.GetAddress(request.OwnerID, installationService.constants.AddressOwners.InstallationRequest)
		customer := installationService.userService.GetUserCredential(request.OwnerID)
		response[i] = installationdto.RequestDetailsResponse{
			ID:           request.ID,
			Name:         request.Name,
			CreatedTime:  request.CreatedAt,
			Status:       request.Status.String(),
			PowerRequest: request.PowerRequest,
			MaxCost:      request.MaxCost,
			BuildingType: request.BuildingType,
			Address:      address,
			Customer:     customer,
		}
	}
	return response
}

func (installationService *InstallationService) AddPanel(panelInfo installationdto.AddPanelRequest) {
	installationService.corporationService.CheckApplicantAccess(panelInfo.CorporationID, panelInfo.OperatorID)
	customer := installationService.userService.FindUserByPhone(panelInfo.CustomerPhone)

	panel := &entity.Panel{
		Name:                 panelInfo.PanelName,
		CorporationID:        panelInfo.CorporationID,
		OperatorID:           panelInfo.OperatorID,
		CustomerID:           customer.ID,
		Power:                panelInfo.Power,
		Area:                 panelInfo.Area,
		BuildingType:         panelInfo.BuildingType,
		Tilt:                 panelInfo.Tilt,
		Azimuth:              panelInfo.Azimuth,
		TotalNumberOfModules: panelInfo.TotalNumberOfModules,
		Address: entity.Address{
			ProvinceID:    panelInfo.Address.ProvinceID,
			CityID:        panelInfo.Address.CityID,
			StreetAddress: panelInfo.Address.StreetAddress,
			PostalCode:    panelInfo.Address.PostalCode,
			HouseNumber:   panelInfo.Address.HouseNumber,
			Unit:          panelInfo.Address.Unit,
		},
	}
	err := installationService.installationRepository.CreatePanel(installationService.db, panel)
	if err != nil {
		panic(err)
	}
}

func (installationService *InstallationService) GetCorporationPanels(listInfo installationdto.PanelListRequest) []installationdto.PanelListResponse {
	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	panels := installationService.installationRepository.FindCorporationPanels(installationService.db, listInfo.CorporationID, paginationModifier, sortingModifier)
	response := make([]installationdto.PanelListResponse, len(panels))
	for i, panel := range panels {
		address := installationService.addressService.GetAddress(panel.ID, installationService.constants.AddressOwners.Panel)
		customer := installationService.userService.GetUserCredential(panel.CustomerID)
		operatior := installationService.userService.GetUserCredential(panel.OperatorID)
		response[i] = installationdto.PanelListResponse{
			ID:                   panel.ID,
			PanelName:            panel.Name,
			CustomerName:         customer.FirstName + " " + customer.LastName,
			CustomerPhone:        customer.Phone,
			Power:                panel.Power,
			Area:                 panel.Area,
			BuildingType:         panel.BuildingType,
			Tilt:                 panel.Tilt,
			Azimuth:              panel.Azimuth,
			TotalNumberOfModules: panel.TotalNumberOfModules,
			Address:              address,
			OperatorName:         operatior.FirstName + " " + operatior.LastName,
		}
	}
	return response

}
