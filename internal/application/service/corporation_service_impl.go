package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type CorporationService struct {
	constants             *bootstrap.Constants
	JWTService            service.JWTService
	db                    database.Database
	CorporationRepository repository.CorporationRepository
	CINService            service.CINService
}

func NewCorporationService(
	constants *bootstrap.Constants,
	jwtService service.JWTService,
	db database.Database,
	corporationRepository repository.CorporationRepository,
	cinService service.CINService,
) *CorporationService {
	return &CorporationService{
		constants:             constants,
		JWTService:            jwtService,
		db:                    db,
		CorporationRepository: corporationRepository,
		CINService:            cinService,
	}
}


func (corporationService *CorporationService) Register(registerInfo corporationdto.RegisterRequest) {
	var conflictErrors exception.ConflictErrors
	_, err := corporationService.CINService.ValidateCIN(registerInfo.CIN)
	if err != nil {
		panic(err)
	}
	corporation, corporationExists := corporationService.CorporationRepository.FindCorporationByCIN(corporationService.db, registerInfo.CIN)
	if corporationExists {
		switch {
		case corporation.Status == corporationService.constants.Status.Approved:
			conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.AlreadyRegistered)
			panic(conflictErrors)
		case corporation.Status == corporationService.constants.Status.Pending :
			conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.Pending)
			panic(conflictErrors)
		}
	}

	corporation = &entity.Corporation{
		Name:     registerInfo.Name,
		CIN:      registerInfo.CIN,
		Password: registerInfo.Password,
		Status:   corporationService.constants.Status.Pending,
	}

	err = corporationService.CorporationRepository.CreateCorporation(corporationService.db, corporation)
	if err != nil {
		panic(err)
	}
}

func (corporationService *CorporationService) Login(loginInfo corporationdto.LoginRequest) corporationdto.CorporationInfoResponse {
	corporation, corporationExist := corporationService.CorporationRepository.FindCorporationByCIN(corporationService.db, loginInfo.CIN)
	var conflictErrors exception.ConflictErrors
	if !corporationExist {
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}
	if corporation.Status != corporationService.constants.Status.Approved {
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}

	if corporation.Password != loginInfo.Password {
		authError := exception.NewInvalidCredentialsError("password not match", nil)
		panic(authError)
	}

	accessToken, refreshToken := corporationService.JWTService.GenerateToken(corporation.ID)
	return corporationdto.CorporationInfoResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Name:         corporation.Name,
	}
}

func (corporationService *CorporationService) GetInstallationRequests(id uint) []corporationdto.InstallationRequestResponse {
	corporation, exist := corporationService.CorporationRepository.FindCorporationByID(corporationService.db, id)
	var conflictErrors exception.ConflictErrors
	switch {
	case !exist:
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	case corporation.Status != corporationService.constants.Status.Approved:
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}

	installationRequests, err := corporationService.CorporationRepository.GetOpenInstallationRequests(corporationService.db, id)
	if err != nil {
		panic(err)
	}


	installationRequestResponses := make([]corporationdto.InstallationRequestResponse, len(installationRequests))
	for i, request := range installationRequests {
		installationRequestResponses[i] = corporationdto.InstallationRequestResponse{
			ID:             request.ID,
			UserID:         request.UserID,
			Area:           request.Area,
			PowerRequested: request.PowerRequested,
			MaxCost:        request.MaxCost,
			Deadline:       request.Deadline,
			BuildingType:   request.BuildingType,
		}
	}
	return installationRequestResponses
}

func (corporationService *CorporationService) SetBid(bidInfo corporationdto.SetBidRequest) {
	corporation, exist := corporationService.CorporationRepository.FindCorporationByID(corporationService.db, bidInfo.CorporationID)
	var conflictErrors exception.ConflictErrors
	switch {
	case !exist:
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	case corporation.Status != corporationService.constants.Status.Approved:
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}

	installationRequest, exist := corporationService.CorporationRepository.FindInstallationRequestByID(corporationService.db, bidInfo.InstallationRequestID)
	switch {
	case !exist:
		conflictErrors.Add(corporationService.constants.Field.InstallationRequest, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	case installationRequest.Status != "open":
		conflictErrors.Add(corporationService.constants.Field.InstallationRequest, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	}
	
	bid := &entity.Bidders{
		RequestType: 			corporationService.constants.RequestType.Installation,
		RequestID: 				bidInfo.InstallationRequestID,
		CorporationID:       	bidInfo.CorporationID,
		MinCost:              	bidInfo.MinCost,
		MaxCost:              	bidInfo.MaxCost,
		MinDeadline:          	bidInfo.MinDeadline,
		MaxDeadline:          	bidInfo.MaxDeadline,
		Description:          	bidInfo.Description,
		InstallationTime:     	bidInfo.InstallationTime,
		Status:              	corporationService.constants.Status.Pending,
	}
	err := corporationService.CorporationRepository.CreateBidder(corporationService.db, bid)
	if err != nil {
		panic(err)
	}
}

func (corporationService *CorporationService) CancelBid(bidInfo corporationdto.CancelBidRequest) {
	corporation, exist := corporationService.CorporationRepository.FindCorporationByID(corporationService.db, bidInfo.CorporationID)
	var conflictErrors exception.ConflictErrors
	switch {
	case !exist:
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	case corporation.Status != corporationService.constants.Status.Approved:
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}


	request, exist := corporationService.CorporationRepository.FindInstallationRequestByID(corporationService.db, bidInfo.InstallationRequestID)
	switch {
	case !exist:
		conflictErrors.Add(corporationService.constants.Field.InstallationRequest, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	case request.Status != "open":
		conflictErrors.Add(corporationService.constants.Field.InstallationRequest, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	}

	bidder, exist := corporationService.CorporationRepository.FindBidderByID(corporationService.db, bidInfo.BidderID)
	switch {
	case !exist:
		conflictErrors.Add(corporationService.constants.Field.Bidder, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	case bidder.CorporationID != bidInfo.CorporationID:
		conflictErrors.Add(corporationService.constants.Field.Bidder, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	case bidder.RequestType != corporationService.constants.RequestType.Installation:
		conflictErrors.Add(corporationService.constants.Field.Bidder, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	case bidder.Status != corporationService.constants.Status.Pending:
		conflictErrors.Add(corporationService.constants.Field.Bidder, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	}
	err := corporationService.CorporationRepository.DeleteBidderByID(corporationService.db, bidInfo.BidderID)
	if err != nil {
		panic(err)
	}
}

func (corporationService *CorporationService) GetBids(corporationID uint) []corporationdto.BidsResponse {
	corporation, exist := corporationService.CorporationRepository.FindCorporationByID(corporationService.db, corporationID)
	var conflictErrors exception.ConflictErrors
	switch {
	case !exist:
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	case corporation.Status != corporationService.constants.Status.Approved:
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}

	bids, err := corporationService.CorporationRepository.GetBids(corporationService.db, corporationID)
	if err != nil {
		panic(err)
	}
	bidResponses := make([]corporationdto.BidsResponse, len(bids))
	for i, bid := range bids {
		bidResponses[i] = corporationdto.BidsResponse{
			ID:                 bid.ID,
			InstallationRequestID:	bid.RequestID,
			MinCost:            	bid.MinCost,
			MaxCost:            	bid.MaxCost,
			MinDeadline:        	bid.MinDeadline,
			MaxDeadline:        	bid.MaxDeadline,
			Description:        	bid.Description,
			InstallationTime:   	bid.InstallationTime,
			Status:             	bid.Status,
		}
	}

	return bidResponses
}
