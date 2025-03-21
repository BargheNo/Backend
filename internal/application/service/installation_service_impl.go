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
	installationRepository repository.InstallationRepository
	db                     database.Database
}

func NewInstallationService(
	constants *bootstrap.Constants,
	addressService service.AddressService,
	installationRepository repository.InstallationRepository,
	db database.Database,
) *InstallationService {
	return &InstallationService{
		constants:              constants,
		addressService:         addressService,
		installationRepository: installationRepository,
		db:                     db,
	}
}

func (installationService *InstallationService) CreateInstallationRequest(requestInfo installationdto.NewInstallationRequest) {
	// get user by id from user service and check complete tag and if not completed -> 403 forbidden
	// compare installed panels names to new request name
	allowedStatus := []enum.InstallationRequestStatus{enum.Active}
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

	address := installationService.addressService.CreateAddress(requestInfo.Address)

	request := &entity.InstallationRequest{
		Name:         requestInfo.Name,
		Status:       enum.Active,
		Area:         requestInfo.Area,
		PowerRequest: requestInfo.Power,
		MaxCost:      requestInfo.MaxCost,
		BuildingType: requestInfo.BuildingType,
		OwnerID:      requestInfo.OwnerID,
		AddressID:    address.ID,
	}
	err := installationService.installationRepository.CreateRequest(installationService.db, request)
	if err != nil {
		panic(err)
	}
}

func (installationService *InstallationService) GetOwnerInstallationRequests(listInfo installationdto.ListOwnerRequestsRequest) []installationdto.ListOwnerRequestsResponse {
	allowedStatus := []enum.InstallationRequestStatus{enum.Active, enum.Cancelled, enum.Expired}
	paginationModifier := repositoryimpl.NewPaginationModifier(listInfo.Limit, listInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
	requests := installationService.installationRepository.FindOwnerRequests(
		installationService.db, listInfo.OwnerID, allowedStatus, paginationModifier, sortingModifier)
	response := make([]installationdto.ListOwnerRequestsResponse, len(requests))
	for i, request := range requests {
		address := installationService.addressService.GetAddress(request.AddressID)
		response[i] = installationdto.ListOwnerRequestsResponse{
			ID:          request.ID,
			Name:        request.Name,
			Status:      request.Status.String(),
			CreatedTime: request.CreatedAt,
			Address:     address,
		}
	}
	return response
}
