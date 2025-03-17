package serviceimpl

import (
	"regexp"

	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enums"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"golang.org/x/crypto/bcrypt"
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

func (corporationService *CorporationService) validatePasswordTests(errors *[]string, test string, password string, tag string) {
	matched, _ := regexp.MatchString(test, password)
	if !matched {
		*errors = append(*errors, tag)
	}
}

func (corporationService *CorporationService) passwordValidation(password string) []string {
	var errors []string

	corporationService.validatePasswordTests(&errors, ".{8,}", password, corporationService.constants.Tag.MinimumLength)
	corporationService.validatePasswordTests(&errors, "[a-z]", password, corporationService.constants.Tag.ContainsLowercase)
	corporationService.validatePasswordTests(&errors, "[A-Z]", password, corporationService.constants.Tag.ContainsUppercase)
	corporationService.validatePasswordTests(&errors, "[0-9]", password, corporationService.constants.Tag.ContainsNumber)
	corporationService.validatePasswordTests(&errors, "[^\\d\\w]", password, corporationService.constants.Tag.ContainsSpecialChar)

	return errors
}

func (corporationService *CorporationService) Register(registerInfo corporationdto.RegisterRequest) {
	var conflictErrors exception.ConflictErrors
	_, err := corporationService.CINService.ValidateCIN(registerInfo.CIN)
	if err != nil {
		panic(err)
	}
	corporation, exist := corporationService.CorporationRepository.FindCorporationByCIN(corporationService.db, registerInfo.CIN)
	if exist && corporation.Status != enums.Rejected.String() {
		conflictErrors.Add(corporationService.constants.Field.CIN, corporationService.constants.Tag.AlreadyRegistered)
		panic(conflictErrors)
	}

	var passwordErrors exception.ValidationErrors
	passwordErrorTags := corporationService.passwordValidation(registerInfo.Password)
	for _, tag := range passwordErrorTags {
		passwordErrors.Add(corporationService.constants.Field.Password, tag)
	}
	if len(passwordErrors.Errors) > 0 {
		panic(passwordErrors)
	}

	hashedPassword, err := hashPassword(registerInfo.Password)
	if err != nil {
		panic(err)
	}

	corporation = &entity.Corporation{
		Name:     registerInfo.Name,
		CIN:      registerInfo.CIN,
		Password: hashedPassword,
		Status:   enums.AwaitingApproval.String(),
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
	if corporation.Status != enums.Approved.String() {
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}

	err := bcrypt.CompareHashAndPassword([]byte(corporation.Password), []byte(loginInfo.Password))
	if err != nil {
		authError := exception.NewInvalidCredentialsError("cin and password not match", nil)
		panic(authError)
	}

	accessToken, refreshToken := corporationService.JWTService.GenerateToken(corporation.ID)
	return corporationdto.CorporationInfoResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Name:         corporation.Name,
	}
}

func (corporationService *CorporationService) GetInstallationRequests(corporationId uint, page int, pageSize int) []corporationdto.InstallationRequestResponse {
	offset := (page - 1) * pageSize
	corporation, exist := corporationService.CorporationRepository.FindCorporationByID(corporationService.db, corporationId)
	var conflictErrors exception.ConflictErrors
	switch {
	case !exist:
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	case corporation.Status != enums.Approved.String():
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}

	installationRequests, err := corporationService.CorporationRepository.GetOpenInstallationRequests(corporationService.db, corporationId, offset, pageSize)
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
			Address:        request.Address,
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
	case corporation.Status != enums.Approved.String():
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}

	installationRequest, exist := corporationService.CorporationRepository.FindInstallationRequestByID(corporationService.db, bidInfo.InstallationRequestID)
	switch {
	case !exist:
		conflictErrors.Add(corporationService.constants.Field.InstallationRequest, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	case installationRequest.Status != enums.Open.String():
		conflictErrors.Add(corporationService.constants.Field.InstallationRequest, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	}

	bid := &entity.Bid{
		RequestType:      enums.InstallationRequest.String(),
		RequestID:        bidInfo.InstallationRequestID,
		CorporationID:    bidInfo.CorporationID,
		MinCost:          bidInfo.MinCost,
		MaxCost:          bidInfo.MaxCost,
		MinDeadline:      bidInfo.MinDeadline,
		MaxDeadline:      bidInfo.MaxDeadline,
		Description:      bidInfo.Description,
		InstallationTime: bidInfo.InstallationTime,
		Status:           enums.Pending.String(),
	}
	err := corporationService.CorporationRepository.CreateBid(corporationService.db, bid)
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
	case corporation.Status != enums.Approved.String():
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}

	request, exist := corporationService.CorporationRepository.FindInstallationRequestByID(corporationService.db, bidInfo.InstallationRequestID)
	switch {
	case !exist:
		conflictErrors.Add(corporationService.constants.Field.InstallationRequest, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	case request.Status != enums.Open.String():
		conflictErrors.Add(corporationService.constants.Field.InstallationRequest, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	}

	bid, exist := corporationService.CorporationRepository.FindBidByID(corporationService.db, bidInfo.BidID)
	switch {
	case !exist:
		conflictErrors.Add(corporationService.constants.Field.Bid, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	case bid.CorporationID != bidInfo.CorporationID:
		conflictErrors.Add(corporationService.constants.Field.Bid, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	case bid.RequestType != enums.InstallationRequest.String():
		conflictErrors.Add(corporationService.constants.Field.Bid, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	case bid.Status != enums.Pending.String():
		conflictErrors.Add(corporationService.constants.Field.Bid, corporationService.constants.Tag.NotExist)
		panic(conflictErrors)
	}
	err := corporationService.CorporationRepository.DeleteBidByID(corporationService.db, bidInfo.BidID)
	if err != nil {
		panic(err)
	}
}

func (corporationService *CorporationService) GetBids(corporationID uint, page int, pageSize int) []corporationdto.BidsResponse {
	offset := (page - 1) * pageSize
	corporation, exist := corporationService.CorporationRepository.FindCorporationByID(corporationService.db, corporationID)
	var conflictErrors exception.ConflictErrors
	switch {
	case !exist:
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	case corporation.Status != enums.Approved.String():
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}

	bids, err := corporationService.CorporationRepository.GetBids(corporationService.db, corporationID, offset, pageSize)
	if err != nil {
		panic(err)
	}
	bidResponses := make([]corporationdto.BidsResponse, len(bids))
	for i, bid := range bids {
		bidResponses[i] = corporationdto.BidsResponse{
			ID:                    bid.ID,
			InstallationRequestID: bid.RequestID,
			MinCost:               bid.MinCost,
			MaxCost:               bid.MaxCost,
			MinDeadline:           bid.MinDeadline,
			MaxDeadline:           bid.MaxDeadline,
			Description:           bid.Description,
			InstallationTime:      bid.InstallationTime,
			Status:                bid.Status,
		}
	}

	return bidResponses
}
