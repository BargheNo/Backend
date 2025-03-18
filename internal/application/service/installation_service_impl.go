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

func (installationService *InstallationService) InstallationRequest(requestInfo installationdto.NewInstallationRequest) {
	// get user by id from user service and check complete tag and if not completed -> 403 forbidden
	allowedStatus := []enum.InstallationRequestStatus{enum.Active}
	inProgressReqs := installationService.installationRepository.FindUserByID(installationService.db, requestInfo.OwnerID, allowedStatus)
	if len(inProgressReqs) >= 5 {
		rateLimitError := exception.NewConcurrentInstallLimitError("", 5, nil)
		panic(rateLimitError)
	}

	address := installationService.addressService.GetAddress(requestInfo.AddressID)

	request := &entity.InstallationRequest{
		Name:         requestInfo.Name,
		Status:       enum.Active,
		Area:         requestInfo.Area,
		PowerRequest: requestInfo.Power,
		MaxCost:      requestInfo.MaxCost,
		BuildingType: requestInfo.BuildingType,
		OwnerID:      requestInfo.OwnerID,
		Address:      *address,
	}
	err := installationService.installationRepository.CreateRequest(installationService.db, request)
	if err != nil {
		panic(err)
	}
}
