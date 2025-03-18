package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enums"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type BidService struct {
	constants             *bootstrap.Constants
	JWTService            service.JWTService
	db                    database.Database
	corporationRepository repository.CorporationRepository
}

func NewBidService(
	constants *bootstrap.Constants,
	jwtService service.JWTService,
	db database.Database,
	corporationRepository repository.CorporationRepository,
) *BidService {
	return &BidService{
		constants:             constants,
		JWTService:            jwtService,
		db:                    db,
		corporationRepository: corporationRepository,
	}
}

func (bidService *BidService) GetInstallationRequests(corporationId uint, page int, pageSize int, sortBy string, ascending bool) []corporationdto.InstallationRequestResponse {
	offset := (page - 1) * pageSize
	dir := "asc"
	if !ascending {
		dir = "desc"
	}
	corporation, exist := bidService.corporationRepository.FindCorporationByID(bidService.db, corporationId)
	var conflictErrors exception.ConflictErrors
	var installationRequests []*entity.InstallationRequest
	var err error
	switch {
	case !exist:
		conflictErrors.Add(bidService.constants.Field.Corporation, bidService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	case corporation.Status != enums.Approved.String():
		conflictErrors.Add(bidService.constants.Field.Corporation, bidService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}

	if sortBy != "" {
		installationRequests, err = bidService.corporationRepository.GetOpenInstallationRequests(bidService.db, corporationId, offset, pageSize, sortBy, dir)
	} else {
		installationRequests, err = bidService.corporationRepository.GetRandomOpenInstallationRequests(bidService.db, corporationId, offset, pageSize)
	}
	
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
			Address:        request.Address.String(),
		}
	}
	return installationRequestResponses
}

func (bidService *BidService) SetBid(bidInfo corporationdto.SetBidRequest) {
	corporation, exist := bidService.corporationRepository.FindCorporationByID(bidService.db, bidInfo.CorporationID)
	var conflictErrors exception.ConflictErrors
	switch {
	case !exist:
		conflictErrors.Add(bidService.constants.Field.Corporation, bidService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	case corporation.Status != enums.Approved.String():
		conflictErrors.Add(bidService.constants.Field.Corporation, bidService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}

	installationRequest, exist := bidService.corporationRepository.FindInstallationRequestByID(bidService.db, bidInfo.InstallationRequestID)
	switch {
	case !exist:
		conflictErrors.Add(bidService.constants.Field.InstallationRequest, bidService.constants.Tag.NotExist)
		panic(conflictErrors)
	case installationRequest.Status != enums.Open.String():
		conflictErrors.Add(bidService.constants.Field.InstallationRequest, bidService.constants.Tag.NotExist)
		panic(conflictErrors)
	}

	bid := &entity.Bid{
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
	err := bidService.corporationRepository.CreateBid(bidService.db, bid)
	if err != nil {
		panic(err)
	}
}

func (bidService *BidService) CancelBid(bidInfo corporationdto.CancelBidRequest) {
	corporation, exist := bidService.corporationRepository.FindCorporationByID(bidService.db, bidInfo.CorporationID)
	var conflictErrors exception.ConflictErrors
	switch {
	case !exist:
		conflictErrors.Add(bidService.constants.Field.Corporation, bidService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	case corporation.Status != enums.Approved.String():
		conflictErrors.Add(bidService.constants.Field.Corporation, bidService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}

	request, exist := bidService.corporationRepository.FindInstallationRequestByID(bidService.db, bidInfo.InstallationRequestID)
	switch {
	case !exist:
		conflictErrors.Add(bidService.constants.Field.InstallationRequest, bidService.constants.Tag.NotExist)
		panic(conflictErrors)
	case request.Status != enums.Open.String():
		conflictErrors.Add(bidService.constants.Field.InstallationRequest, bidService.constants.Tag.NotExist)
		panic(conflictErrors)
	}

	bid, exist := bidService.corporationRepository.FindBidByID(bidService.db, bidInfo.BidID)
	switch {
	case !exist:
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.NotExist)
		panic(conflictErrors)
	case bid.CorporationID != bidInfo.CorporationID:
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.NotExist)
		panic(conflictErrors)
	case bid.Status != enums.Pending.String():
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.NotExist)
		panic(conflictErrors)
	}
	err := bidService.corporationRepository.DeleteBidByID(bidService.db, bidInfo.BidID)
	if err != nil {
		panic(err)
	}
}

func (bidService *BidService) GetBids(corporationID uint, page int, pageSize int, sortBy string, ascending bool) []corporationdto.BidsResponse {
	offset := (page - 1) * pageSize
	dir := "asc"
	if !ascending {
		dir = "desc"
	}
	corporation, exist := bidService.corporationRepository.FindCorporationByID(bidService.db, corporationID)
	var conflictErrors exception.ConflictErrors
	switch {
	case !exist:
		conflictErrors.Add(bidService.constants.Field.Corporation, bidService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	case corporation.Status != enums.Approved.String():
		conflictErrors.Add(bidService.constants.Field.Corporation, bidService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}
	if sortBy == "" {
		sortBy = "id"
	}
	
	bids, err := bidService.corporationRepository.GetBids(bidService.db, corporationID, offset, pageSize, sortBy, dir)
	
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
