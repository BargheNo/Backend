package serviceimpl

import (
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	chatdto "github.com/BargheNo/Backend/internal/application/dto/chat"
	guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"
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
	guaranteeService       service.GuaranteeService
	chatService            service.ChatService
	installationRepository repository.InstallationRepository
	db                     database.Database
}

type InstallationServiceDeps struct {
	Constants              *bootstrap.Constants
	AddressService         service.AddressService
	UserService            service.UserService
	CorporationService     service.CorporationService
	GuaranteeService       service.GuaranteeService
	ChatService            service.ChatService
	InstallationRepository repository.InstallationRepository
	DB                     database.Database
}

func NewInstallationService(deps InstallationServiceDeps) *InstallationService {
	return &InstallationService{
		constants:              deps.Constants,
		addressService:         deps.AddressService,
		userService:            deps.UserService,
		corporationService:     deps.CorporationService,
		guaranteeService:       deps.GuaranteeService,
		chatService:            deps.ChatService,
		installationRepository: deps.InstallationRepository,
		db:                     deps.DB,
	}
}

func (installationService *InstallationService) GetRequestStatuses() []installationdto.EnumStatusResponse {
	statuses := enum.GetAllInstallationRequestStatuses()
	response := make([]installationdto.EnumStatusResponse, len(statuses))
	for i, status := range statuses {
		response[i] = installationdto.EnumStatusResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response
}

func (installationService *InstallationService) GetBuildingTypes() []installationdto.EnumStatusResponse {
	types := enum.GetAllBuildingTypes()
	response := make([]installationdto.EnumStatusResponse, len(types))
	for i, buildingType := range types {
		response[i] = installationdto.EnumStatusResponse{
			ID:   uint(buildingType),
			Name: buildingType.String(),
		}
	}
	return response
}

func (installationService *InstallationService) ValidateRequestOwnership(requestID, ownerID uint) (installationdto.PublicRequestDetailsResponse, error) {
	request, exist := installationService.installationRepository.FindRequestByOwner(installationService.db, requestID, ownerID)
	if !exist {
		return installationdto.PublicRequestDetailsResponse{}, exception.NotFoundError{Item: installationService.constants.Field.InstallationRequest}
	}

	customer := installationService.userService.GetUserCredential(request.OwnerID)
	address := installationService.addressService.GetAddress(request.ID, installationService.constants.AddressOwners.InstallationRequest)

	return installationdto.PublicRequestDetailsResponse{
		ID:           request.ID,
		Name:         request.Name,
		Status:       request.Status.String(),
		PowerRequest: request.PowerRequest,
		Description:  request.Description,
		BuildingType: request.BuildingType.String(),
		Area:         request.Area,
		MaxCost:      request.MaxCost,
		Customer:     customer,
		Address:      address,
	}, nil
}

func (installationService *InstallationService) CreateInstallationRequest(request installationdto.NewInstallationRequest) {
	// compare installed panels names to new request name
	if ok := installationService.userService.IsUserActive(request.OwnerID); !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: installationService.constants.Field.InstallationRequest,
		}
		panic(forbiddenError)
	}

	allowedStatus := []enum.InstallationRequestStatus{enum.InstallationRequestStatusActive}
	_, exist := installationService.installationRepository.FindOwnerRequestByName(installationService.db, request.OwnerID, allowedStatus, request.Name)
	if exist {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(installationService.constants.Field.Name, installationService.constants.Tag.AlreadyRegistered)
		panic(conflictErrors)
	}
	inProgressReqs := installationService.installationRepository.FindOwnerRequests(installationService.db, request.OwnerID, allowedStatus)
	if len(inProgressReqs) >= 5 {
		rateLimitError := exception.NewConcurrentInstallLimitError("", 5, nil)
		panic(rateLimitError)
	}

	installationRequest := &entity.InstallationRequest{
		Name:         request.Name,
		Status:       enum.InstallationRequestStatusActive,
		Area:         request.Area,
		PowerRequest: request.Power,
		MaxCost:      request.MaxCost,
		BuildingType: enum.BuildingType(request.BuildingType),
		OwnerID:      request.OwnerID,
		Description:  request.Description,
		Address: entity.Address{
			ProvinceID:    request.Address.ProvinceID,
			CityID:        request.Address.CityID,
			StreetAddress: request.Address.StreetAddress,
			PostalCode:    request.Address.PostalCode,
			HouseNumber:   request.Address.HouseNumber,
			Unit:          request.Address.Unit,
		},
	}
	err := installationService.installationRepository.CreateRequest(installationService.db, installationRequest)
	if err != nil {
		panic(err)
	}
}

func (installationService *InstallationService) GetOwnerInstallationRequests(request installationdto.CustomerRequestsListRequest) []installationdto.AnonymousRequestsResponse {
	status := enum.InstallationRequestStatus(request.Status)
	allowedStatus := []enum.InstallationRequestStatus{status}
	if status == enum.InstallationRequestStatusAll {
		allowedStatus = enum.GetAllInstallationRequestStatuses()
	}

	paginationModifier := repositoryimpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	requests := installationService.installationRepository.FindOwnerRequests(
		installationService.db, request.OwnerID, allowedStatus, paginationModifier, sortingModifier)
	response := make([]installationdto.AnonymousRequestsResponse, len(requests))

	for i, request := range requests {
		address := installationService.addressService.GetAddress(request.ID, installationService.constants.AddressOwners.InstallationRequest)
		response[i] = installationdto.AnonymousRequestsResponse{
			ID:           request.ID,
			Name:         request.Name,
			CreatedTime:  request.CreatedAt,
			Status:       request.Status.String(),
			PowerRequest: request.PowerRequest,
			MaxCost:      request.MaxCost,
			BuildingType: request.BuildingType.String(),
			Address:      address,
		}
	}
	return response
}

func (installationService *InstallationService) GetOwnerInstallationRequest(request installationdto.GetOwnerRequest) installationdto.AnonymousRequestsResponse {
	installationRequest, exist := installationService.installationRepository.FindRequestByOwner(installationService.db, request.InstallationID, request.OwnerID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.InstallationRequest}
		panic(notFoundError)
	}

	address := installationService.addressService.GetAddress(request.InstallationID, installationService.constants.AddressOwners.InstallationRequest)
	return installationdto.AnonymousRequestsResponse{
		ID:           installationRequest.ID,
		Name:         installationRequest.Name,
		CreatedTime:  installationRequest.CreatedAt,
		Status:       installationRequest.Status.String(),
		PowerRequest: installationRequest.PowerRequest,
		MaxCost:      installationRequest.MaxCost,
		BuildingType: installationRequest.BuildingType.String(),
		Address:      address,
	}
}

func (installationService *InstallationService) ChangeInstallationRequestStatus(request installationdto.ChangeRequestStatusRequest) {
	installationRequest, exist := installationService.installationRepository.FindRequestByOwner(installationService.db, request.RequestID, request.OwnerID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.InstallationRequest}
		panic(notFoundError)
	}

	if installationRequest.Status == enum.InstallationRequestStatusCancelled {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(installationService.constants.Field.InstallationRequest, installationService.constants.Tag.AlreadyCanceled)
		panic(conflictErrors)
	}

	installationRequest.Status = request.Status
	if err := installationService.installationRepository.UpdateRequest(installationService.db, installationRequest); err != nil {
		panic(err)
	}
}

func (installationService *InstallationService) GetAnonymousInstallationRequests(request installationdto.CorporationPanelListRequest) []installationdto.AnonymousRequestsResponse {
	installationService.corporationService.CheckApplicantAccess(request.CorporationID, request.OperatorID)

	allowedStatus := []enum.InstallationRequestStatus{enum.InstallationRequestStatusActive}

	paginationModifier := repositoryimpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	installationRequests := installationService.installationRepository.FindRequestsByStatus(installationService.db, allowedStatus, paginationModifier, sortingModifier)
	response := make([]installationdto.AnonymousRequestsResponse, len(installationRequests))

	for i, installationRequest := range installationRequests {
		address := installationService.addressService.GetAddress(installationRequest.ID, installationService.constants.AddressOwners.InstallationRequest)
		response[i] = installationdto.AnonymousRequestsResponse{
			ID:           installationRequest.ID,
			Name:         installationRequest.Name,
			CreatedTime:  installationRequest.CreatedAt,
			Status:       installationRequest.Status.String(),
			PowerRequest: installationRequest.PowerRequest,
			MaxCost:      installationRequest.MaxCost,
			BuildingType: installationRequest.BuildingType.String(),
			Address:      address,
		}
	}
	return response
}

func (installationService *InstallationService) GetAnonymousInstallationRequest(request installationdto.CorporationPanelRequest) installationdto.AnonymousRequestsResponse {
	installationService.corporationService.CheckApplicantAccess(request.CorporationID, request.OperatorID)

	installationRequest, exist := installationService.installationRepository.FindRequestByID(installationService.db, request.InstallationID)
	if !exist || installationRequest.Status != enum.InstallationRequestStatusActive {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.InstallationRequest}
		panic(notFoundError)
	}

	address := installationService.addressService.GetAddress(installationRequest.ID, installationService.constants.AddressOwners.InstallationRequest)
	return installationdto.AnonymousRequestsResponse{
		ID:           installationRequest.ID,
		Name:         installationRequest.Name,
		CreatedTime:  installationRequest.CreatedAt,
		Status:       installationRequest.Status.String(),
		PowerRequest: installationRequest.PowerRequest,
		MaxCost:      installationRequest.MaxCost,
		BuildingType: installationRequest.BuildingType.String(),
		Address:      address,
	}
}

func (installationService *InstallationService) GetPublicInstallationRequest(requestID uint) installationdto.PublicRequestDetailsResponse {
	installationRequest, exist := installationService.installationRepository.FindRequestByID(installationService.db, requestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.InstallationRequest}
		panic(notFoundError)
	}
	customer := installationService.userService.GetUserCredential(installationRequest.OwnerID)
	address := installationService.addressService.GetAddress(installationRequest.ID, installationService.constants.AddressOwners.InstallationRequest)
	return installationdto.PublicRequestDetailsResponse{
		ID:           installationRequest.ID,
		Name:         installationRequest.Name,
		Status:       installationRequest.Status.String(),
		PowerRequest: installationRequest.PowerRequest,
		Description:  installationRequest.Description,
		BuildingType: installationRequest.BuildingType.String(),
		Area:         installationRequest.Area,
		MaxCost:      installationRequest.MaxCost,
		Customer:     customer,
		Address:      address,
	}
}

func (installationService *InstallationService) GetInstallationRequestsByAdmin(request installationdto.AdminRequestsListRequest) []installationdto.PublicRequestDetailsResponse {
	allowedStatuses := []enum.InstallationRequestStatus{enum.InstallationRequestStatus(request.Status)}
	if enum.InstallationRequestStatus(request.Status) == enum.InstallationRequestStatusAll {
		allowedStatuses = enum.GetAllInstallationRequestStatuses()
	}

	paginationModifier := repositoryimpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	installationRequests := installationService.installationRepository.FindRequestsByStatus(installationService.db, allowedStatuses, sortingModifier, paginationModifier)
	response := make([]installationdto.PublicRequestDetailsResponse, len(installationRequests))

	for i, installationRequest := range installationRequests {
		customer := installationService.userService.GetUserCredential(installationRequest.OwnerID)
		address := installationService.addressService.GetAddress(installationRequest.ID, installationService.constants.AddressOwners.InstallationRequest)

		response[i] = installationdto.PublicRequestDetailsResponse{
			ID:           installationRequest.ID,
			Name:         installationRequest.Name,
			Status:       installationRequest.Status.String(),
			PowerRequest: installationRequest.PowerRequest,
			Description:  installationRequest.Description,
			BuildingType: installationRequest.BuildingType.String(),
			Area:         installationRequest.Area,
			MaxCost:      installationRequest.MaxCost,
			Customer:     customer,
			Address:      address,
		}
	}

	return response
}

func (installationService *InstallationService) CompleteInstallationRequest(request installationdto.CompleteBidInstallationRequest) {
	installationService.corporationService.CheckApplicantAccess(request.CorporationID, request.OperatorID)

	panel, exist := installationService.installationRepository.FindCorporationPanel(installationService.db, request.PanelID, request.CorporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.Panel}
		panic(notFoundError)
	}

	panel.Tilt = request.Tilt
	panel.Azimuth = request.Azimuth
	panel.TotalNumberOfModules = request.NumberOfModules

	if err := installationService.installationRepository.UpdatePanel(installationService.db, panel); err != nil {
		panic(err)
	}
}

func (installationService *InstallationService) UpdateInstallationRequestByAdmin(newRequest installationdto.UpdateInstallationRequest) {
	installationRequest, exist := installationService.installationRepository.FindRequestByID(installationService.db, newRequest.RequestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.InstallationRequest}
		panic(notFoundError)
	}

	if newRequest.Name != nil {
		installationRequest.Name = *newRequest.Name
	}
	if newRequest.Area != nil {
		installationRequest.Area = *newRequest.Area
	}
	if newRequest.Power != nil {
		installationRequest.PowerRequest = *newRequest.Power
	}
	if newRequest.MaxCost != nil {
		installationRequest.MaxCost = *newRequest.MaxCost
	}
	if newRequest.BuildingType != nil {
		installationRequest.BuildingType = enum.BuildingType(*newRequest.BuildingType)
	}
	if newRequest.Status != nil {
		installationRequest.Status = enum.InstallationRequestStatus(*newRequest.Status)
	}
	if newRequest.Description != nil {
		installationRequest.Description = *newRequest.Description
	}

	if err := installationService.installationRepository.UpdateRequest(installationService.db, installationRequest); err != nil {
		panic(err)
	}
}

func (installationService *InstallationService) DeleteInstallationRequest(requestID uint) {
	installationRequest, exist := installationService.installationRepository.FindRequestByID(installationService.db, requestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.InstallationRequest}
		panic(notFoundError)
	}

	if err := installationService.installationRepository.DeleteRequest(installationService.db, installationRequest); err != nil {
		panic(err)
	}
}

func (installationService *InstallationService) ValidatePanelOwnership(panelID, userID uint) error {
	_, exist := installationService.installationRepository.FindPanelByOwner(installationService.db, panelID, userID)
	if !exist {
		return exception.NotFoundError{Item: installationService.constants.Field.InstallationRequest}
	}
	return nil
}

func (installationService *InstallationService) ValidatePanelGuarantee(panelID uint) error {
	panel, exist := installationService.installationRepository.FindPanelByID(installationService.db, panelID)
	if !exist {
		return exception.NotFoundError{Item: installationService.constants.Field.InstallationRequest}
	}
	if panel.GuaranteeStatus != enum.PanelGuaranteeStatusActive {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(installationService.constants.Field.Guarantee, installationService.constants.Tag.NotActive)
		return conflictErrors
	}
	return nil
}

// TODO: nor done remain for the bid/bidID/request && maybe remove:FindPanelByNameAndCustomerID
func (installationService *InstallationService) AddPanel(panelInfo installationdto.AddPanelRequest) {
	installationService.corporationService.CheckApplicantAccess(panelInfo.CorporationID, panelInfo.OperatorID)
	if ok := installationService.userService.IsUserActive(panelInfo.OperatorID); !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: installationService.constants.Field.Panel,
		}
		panic(forbiddenError)
	}

	customer := installationService.userService.FindUserByPhone(panelInfo.CustomerPhone)

	if _, exist := installationService.installationRepository.FindPanelByNameAndCustomerID(installationService.db, panelInfo.Name, customer.ID); exist {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(installationService.constants.Field.Name, installationService.constants.Tag.AlreadyExist)
		panic(conflictErrors)
	}

	panelGuaranteeStatus := enum.PanelGuaranteeStatusEmpty
	if panelInfo.GuaranteeID != nil {
		if err := installationService.guaranteeService.ValidateActiveGuaranteeOwnerShip(*panelInfo.GuaranteeID, panelInfo.CorporationID); err != nil {
			panic(err)
		}
		panelGuaranteeStatus = enum.PanelGuaranteeStatusActive
	}

	panel := &entity.Panel{
		Name:                 panelInfo.Name,
		Status:               panelInfo.Status,
		BuildingType:         enum.BuildingType(panelInfo.BuildingType),
		Power:                panelInfo.Power,
		Area:                 panelInfo.Area,
		Tilt:                 panelInfo.Tilt,
		Azimuth:              panelInfo.Azimuth,
		TotalNumberOfModules: panelInfo.TotalNumberOfModules,
		CorporationID:        panelInfo.CorporationID,
		OperatorID:           panelInfo.OperatorID,
		CustomerID:           customer.ID,
		GuaranteeStatus:      panelGuaranteeStatus,
		GuaranteeID:          panelInfo.GuaranteeID,
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

	request := chatdto.CreateOrGetUserRoomRequest{
		CorporationID: panel.CorporationID,
		UserID:        customer.ID,
	}
	installationService.chatService.CreateOrGetRoom(request)
}

func (installationService *InstallationService) GetCorporationPanels(listInfo installationdto.CorporationPanelListRequest) []installationdto.CorporationPanelListResponse {
	installationService.corporationService.CheckApplicantAccess(listInfo.CorporationID, listInfo.OperatorID)

	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	allowedStatus := []enum.PanelStatus{enum.PanelStatus(listInfo.Status)}
	if enum.PanelStatus(listInfo.Status) == enum.PanelStatusAll {
		allowedStatus = enum.GetAllPanelStatuses()
	}

	panels := installationService.installationRepository.FindCorporationPanels(installationService.db, listInfo.CorporationID, allowedStatus, paginationModifier, sortingModifier)
	response := make([]installationdto.CorporationPanelListResponse, len(panels))

	for i, panel := range panels {
		address := installationService.addressService.GetAddress(panel.ID, installationService.constants.AddressOwners.Panel)
		customer := installationService.userService.GetUserCredential(panel.CustomerID)
		operator := installationService.userService.GetUserCredential(panel.OperatorID)

		response[i] = installationdto.CorporationPanelListResponse{
			ID:                   panel.ID,
			Name:                 panel.Name,
			Status:               panel.Status.String(),
			BuildingType:         panel.BuildingType.String(),
			Area:                 panel.Area,
			Power:                panel.Power,
			Tilt:                 panel.Tilt,
			Azimuth:              panel.Azimuth,
			TotalNumberOfModules: panel.TotalNumberOfModules,
			GuaranteeStatus:      panel.GuaranteeStatus.String(),
			Operator:             operator,
			Customer:             customer,
			Address:              address,
		}
	}
	return response
}

func (installationService *InstallationService) GetCorporationPanel(request installationdto.CorporationPanelRequest) installationdto.CorporationPanelResponse {
	installationService.corporationService.CheckApplicantAccess(request.CorporationID, request.OperatorID)

	panel, exist := installationService.installationRepository.FindCorporationPanel(installationService.db, request.InstallationID, request.CorporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.Panel}
		panic(notFoundError)
	}

	address := installationService.addressService.GetAddress(panel.ID, installationService.constants.AddressOwners.Panel)
	customer := installationService.userService.GetUserCredential(panel.CustomerID)
	operator := installationService.userService.GetUserCredential(panel.OperatorID)

	var guarantee guaranteedto.GuaranteeResponse
	var err error
	if panel.GuaranteeID != nil {
		guarantee, err = installationService.guaranteeService.GetGuarantee(*panel.GuaranteeID)
		if err != nil {
			panic(err)
		}
	}

	return installationdto.CorporationPanelResponse{
		ID:                   panel.ID,
		Name:                 panel.Name,
		Status:               panel.Status.String(),
		BuildingType:         panel.BuildingType.String(),
		Area:                 panel.Area,
		Power:                panel.Power,
		Tilt:                 panel.Tilt,
		Azimuth:              panel.Azimuth,
		TotalNumberOfModules: panel.TotalNumberOfModules,
		GuaranteeStatus:      panel.GuaranteeStatus.String(),
		Operator:             operator,
		Customer:             customer,
		Address:              address,
		Guarantee:            guarantee,
	}
}

func (installationService *InstallationService) GetCustomerPanels(listInfo installationdto.CustomerPanelListRequest) []installationdto.CustomerPanelListResponse {
	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	allowedStatus := []enum.PanelStatus{enum.PanelStatus(listInfo.Status)}
	if enum.PanelStatus(listInfo.Status) == enum.PanelStatusAll {
		allowedStatus = enum.GetAllPanelStatuses()
	}

	panels := installationService.installationRepository.FindCustomerPanels(installationService.db, listInfo.OwnerID, allowedStatus, paginationModifier, sortingModifier)
	response := make([]installationdto.CustomerPanelListResponse, len(panels))

	for i, panel := range panels {
		address := installationService.addressService.GetAddress(panel.ID, installationService.constants.AddressOwners.Panel)
		corporation := installationService.corporationService.GetCorporationCredentials(panel.CorporationID)

		response[i] = installationdto.CustomerPanelListResponse{
			ID:                   panel.ID,
			Name:                 panel.Name,
			Status:               panel.Status.String(),
			BuildingType:         panel.BuildingType.String(),
			Area:                 panel.Area,
			Power:                panel.Power,
			Tilt:                 panel.Tilt,
			Azimuth:              panel.Azimuth,
			TotalNumberOfModules: panel.TotalNumberOfModules,
			GuaranteeStatus:      panel.GuaranteeStatus.String(),
			Corporation:          corporation,
			Address:              address,
		}
	}
	return response
}

func (installationService *InstallationService) GetCustomerPanel(panelInfo installationdto.GetOwnerRequest) installationdto.CustomerPanelResponse {
	panel, exist := installationService.installationRepository.FindCustomerPanel(installationService.db, panelInfo.InstallationID, panelInfo.OwnerID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.Panel}
		panic(notFoundError)
	}

	address := installationService.addressService.GetAddress(panel.ID, installationService.constants.AddressOwners.Panel)
	corporation := installationService.corporationService.GetCorporationCredentials(panel.CorporationID)

	var guarantee guaranteedto.GuaranteeResponse
	var err error
	if panel.GuaranteeID != nil {
		guarantee, err = installationService.guaranteeService.GetGuarantee(*panel.GuaranteeID)
		if err != nil {
			panic(err)
		}
	}

	return installationdto.CustomerPanelResponse{
		ID:                   panel.ID,
		Name:                 panel.Name,
		Status:               panel.Status.String(),
		BuildingType:         panel.BuildingType.String(),
		Area:                 panel.Area,
		Power:                panel.Power,
		Tilt:                 panel.Tilt,
		Azimuth:              panel.Azimuth,
		TotalNumberOfModules: panel.TotalNumberOfModules,
		GuaranteeStatus:      panel.GuaranteeStatus.String(),
		Corporation:          corporation,
		Address:              address,
		Guarantee:            guarantee,
	}
}

func (installationService *InstallationService) GetPanelByAdmin(panelID uint) installationdto.AdminPanelResponse {
	panel, exist := installationService.installationRepository.FindPanelByID(installationService.db, panelID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.Panel}
		panic(notFoundError)
	}

	customer := installationService.userService.GetUserCredential(panel.CustomerID)
	operator := installationService.userService.GetUserCredential(panel.OperatorID)
	corporation := installationService.corporationService.GetCorporationCredentials(panel.CorporationID)
	address := installationService.addressService.GetAddress(panel.ID, installationService.constants.AddressOwners.Panel)

	var guarantee guaranteedto.GuaranteeResponse
	if panel.Guarantee != nil {
		guarantee, _ = installationService.guaranteeService.GetGuarantee(*panel.GuaranteeID)
	}

	return installationdto.AdminPanelResponse{
		ID:                   panelID,
		Name:                 panel.Name,
		Status:               panel.Status.String(),
		BuildingType:         panel.BuildingType.String(),
		Area:                 panel.Area,
		Power:                panel.Power,
		Tilt:                 panel.Tilt,
		Azimuth:              panel.Azimuth,
		TotalNumberOfModules: panel.TotalNumberOfModules,
		Operator:             operator,
		Customer:             customer,
		Corporation:          corporation,
		Address:              address,
		Guarantee:            guarantee,
	}
}

func (installationService *InstallationService) ViolatePanelGuaranteeStatus(request installationdto.CreateViolatePanelGuaranteeRequest) uint {
	installationService.corporationService.CheckApplicantAccess(request.CorporationID, request.OperatorID)

	panel, exist := installationService.installationRepository.FindCorporationPanel(installationService.db, request.PanelID, request.CorporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.Panel}
		panic(notFoundError)
	}

	if panel.GuaranteeStatus == enum.PanelGuaranteeStatusEmpty {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.Guarantee}
		panic(notFoundError)
	}

	if panel.GuaranteeStatus != enum.PanelGuaranteeStatusActive {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(installationService.constants.Field.Guarantee, installationService.constants.Tag.NotActive)
		panic(conflictErrors)
	}

	violationID := installationService.guaranteeService.CreateGuaranteeViolation(request.GuaranteeViolation)
	panel.GuaranteeStatus = enum.PanelGuaranteeStatusVoided

	if err := installationService.installationRepository.UpdatePanel(installationService.db, panel); err != nil {
		panic(err)
	}

	return violationID
}

func (installationService *InstallationService) ClearPanelGuaranteeViolation(violationInfo installationdto.GetCorporationGuaranteeViolationRequest) {
	installationService.corporationService.CheckApplicantAccess(violationInfo.CorporationID, violationInfo.OperatorID)

	panel, exist := installationService.installationRepository.FindCorporationPanel(installationService.db, violationInfo.PanelID, violationInfo.CorporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.Panel}
		panic(notFoundError)
	}

	if panel.GuaranteeStatus == enum.PanelGuaranteeStatusEmpty {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.Guarantee}
		panic(notFoundError)
	}

	if panel.GuaranteeStatus != enum.PanelGuaranteeStatusVoided {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.GuaranteeViolation}
		panic(notFoundError)
	}

	installationService.guaranteeService.RemovePanelGuaranteeViolation(violationInfo.PanelID)

	if panel.GuaranteeEndDate.Before(time.Now()) {
		panel.GuaranteeStatus = enum.PanelGuaranteeStatusExpired
	} else {
		panel.GuaranteeStatus = enum.PanelGuaranteeStatusActive
	}

	if err := installationService.installationRepository.UpdatePanel(installationService.db, panel); err != nil {
		panic(err)
	}
}

func (installationService *InstallationService) GetCorporationPanelGuaranteeViolation(violationInfo installationdto.GetCorporationGuaranteeViolationRequest) (guaranteedto.CorporationGuaranteeViolationResponse, error) {
	var violation guaranteedto.CorporationGuaranteeViolationResponse

	installationService.corporationService.CheckApplicantAccess(violationInfo.CorporationID, violationInfo.OperatorID)

	panel, exist := installationService.installationRepository.FindCorporationPanel(installationService.db, violationInfo.PanelID, violationInfo.CorporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.Panel}
		return violation, notFoundError
	}

	if panel.GuaranteeStatus == enum.PanelGuaranteeStatusEmpty {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.Guarantee}
		return violation, notFoundError
	}

	if panel.GuaranteeStatus != enum.PanelGuaranteeStatusVoided {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.GuaranteeViolation}
		return violation, notFoundError
	}

	violation, err := installationService.guaranteeService.GetCorporationPanelGuaranteeViolation(violationInfo.PanelID)
	if err != nil {
		panic(err)
	}

	return violation, nil
}

func (installationService *InstallationService) GetCustomerPanelGuaranteeViolation(violationInfo installationdto.GetCustomerGuaranteeViolationRequest) (guaranteedto.CustomerGuaranteeViolationResponse, error) {
	var violation guaranteedto.CustomerGuaranteeViolationResponse

	panel, exist := installationService.installationRepository.FindCustomerPanel(installationService.db, violationInfo.PanelID, violationInfo.OwnerID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.Panel}
		return violation, notFoundError
	}

	if panel.GuaranteeStatus == enum.PanelGuaranteeStatusEmpty {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.Guarantee}
		return violation, notFoundError
	}

	if panel.GuaranteeStatus != enum.PanelGuaranteeStatusVoided {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.GuaranteeViolation}
		return violation, notFoundError
	}

	violation, err := installationService.guaranteeService.GetCustomerPanelGuaranteeViolation(violationInfo.PanelID)
	if err != nil {
		panic(err)
	}

	return violation, nil
}

func (installationService *InstallationService) UpdatePanelGuaranteeViolation(violationInfo installationdto.UpdateGuaranteeViolationRequest) {
	installationService.corporationService.CheckApplicantAccess(violationInfo.CorporationID, violationInfo.OperatorID)

	panel, exist := installationService.installationRepository.FindCorporationPanel(installationService.db, violationInfo.PanelID, violationInfo.CorporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.Panel}
		panic(notFoundError)
	}

	if panel.GuaranteeStatus == enum.PanelGuaranteeStatusEmpty {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.Guarantee}
		panic(notFoundError)
	}

	if panel.GuaranteeStatus != enum.PanelGuaranteeStatusVoided {
		notFoundError := exception.NotFoundError{Item: installationService.constants.Field.GuaranteeViolation}
		panic(notFoundError)
	}

	request := guaranteedto.UpdateGuaranteeViolationRequest{
		PanelID:    violationInfo.PanelID,
		OperatorID: violationInfo.OperatorID,
		Reason:     violationInfo.Reason,
		Details:    violationInfo.Details,
	}
	installationService.guaranteeService.UpdateGuaranteeViolation(request)
}
