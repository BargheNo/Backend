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
	bidRepository		  repository.BidRepository
	corporationService   service.CorporationService
}

func NewBidService(
	constants *bootstrap.Constants,
	jwtService service.JWTService,
	db database.Database,
	bidRepository repository.BidRepository,
	corporationService service.CorporationService,
) *BidService {
	return &BidService{
		constants:             constants,
		JWTService:            jwtService,
		db:                    db,
		bidRepository:         bidRepository,
		corporationService:    corporationService,
	}
}



func (bidService *BidService) GetInstallationRequests(corporationId uint, page int, pageSize int, sortBy string, ascending bool) []corporationdto.InstallationRequestResponse {
	offset := (page - 1) * pageSize
	dir := "asc"
	if !ascending {
		dir = "desc"
	}
	_, exist := bidService.corporationService.GetCorporationByID(corporationId)
	if !exist {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(bidService.constants.Field.Corporation, bidService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}
	var installationRequests []*entity.InstallationRequest
	var err error
	

	if sortBy != "" {
		installationRequests, err = bidService.bidRepository.GetOpenInstallationRequests(bidService.db, corporationId, offset, pageSize, sortBy, dir)
	} else {
		installationRequests, err = bidService.bidRepository.GetRandomOpenInstallationRequests(bidService.db, corporationId, offset, pageSize)
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
	_, exist := bidService.corporationService.GetCorporationByID(bidInfo.CorporationID)
	var conflictErrors exception.ConflictErrors
	if !exist {
		conflictErrors.Add(bidService.constants.Field.Corporation, bidService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}
	installationRequest, exist := bidService.bidRepository.FindInstallationRequestByID(bidService.db, bidInfo.InstallationRequestID)
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
	err := bidService.bidRepository.CreateBid(bidService.db, bid)
	if err != nil {
		panic(err)
	}
}

func (bidService *BidService) CancelBid(bidInfo corporationdto.CancelBidRequest) {
	_, exist := bidService.corporationService.GetCorporationByID(bidInfo.CorporationID)
	var conflictErrors exception.ConflictErrors
	if !exist {
		conflictErrors.Add(bidService.constants.Field.Corporation, bidService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}

	request, exist := bidService.bidRepository.FindInstallationRequestByID(bidService.db, bidInfo.InstallationRequestID)
	switch {
	case !exist:
		conflictErrors.Add(bidService.constants.Field.InstallationRequest, bidService.constants.Tag.NotExist)
		panic(conflictErrors)
	case request.Status != enums.Open.String():
		conflictErrors.Add(bidService.constants.Field.InstallationRequest, bidService.constants.Tag.NotExist)
		panic(conflictErrors)
	}

	bid, exist := bidService.bidRepository.FindBidByID(bidService.db, bidInfo.BidID)
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
	err := bidService.bidRepository.DeleteBidByID(bidService.db, bidInfo.BidID)
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
	_, exist := bidService.corporationService.GetCorporationByID(corporationID)
	var conflictErrors exception.ConflictErrors
	if !exist {
		conflictErrors.Add(bidService.constants.Field.Corporation, bidService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}
	if sortBy == "" {
		sortBy = "id"
	}
	
	bids, err := bidService.bidRepository.GetBids(bidService.db, corporationID, offset, pageSize, sortBy, dir)
	
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
