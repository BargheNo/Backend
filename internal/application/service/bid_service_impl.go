package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enums"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type BidService struct {
	constants          *bootstrap.Constants
	JWTService         service.JWTService
	db                 database.Database
	bidRepository      repository.BidRepository
	corporationService service.CorporationService
}

func NewBidService(
	constants *bootstrap.Constants,
	jwtService service.JWTService,
	db database.Database,
	bidRepository repository.BidRepository,
	corporationService service.CorporationService,
) *BidService {
	return &BidService{
		constants:          constants,
		JWTService:         jwtService,
		db:                 db,
		bidRepository:      bidRepository,
		corporationService: corporationService,
	}
}

func (bidService *BidService) GetInstallationRequests(corporationId uint, page int, pageSize int, sortBy string, ascending bool) []biddto.InstallationRequestResponse {
	offset := (page - 1) * pageSize
	dir := "ASC"
	if !ascending {
		dir = "DESC"
	}
	_, exist := bidService.corporationService.GetCorporationByID(corporationId)
	if !exist {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(bidService.constants.Field.Corporation, bidService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}
	var installationRequests []*entity.InstallationRequest

	if sortBy != "" {
		installationRequests = bidService.bidRepository.GetOpenInstallationRequests(bidService.db, corporationId, offset, pageSize, sortBy, dir)
	} else {
		installationRequests = bidService.bidRepository.GetRandomOpenInstallationRequests(bidService.db, corporationId, offset, pageSize)
	}
	installationRequestResponses := make([]biddto.InstallationRequestResponse, len(installationRequests))
	for i, request := range installationRequests {
		installationRequestResponses[i] = biddto.InstallationRequestResponse{
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

func (bidService *BidService) SetBid(bidInfo biddto.SetBidRequest) {
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
	case installationRequest.Status != enums.Open:
		conflictErrors.Add(bidService.constants.Field.InstallationRequest, bidService.constants.Tag.NotExist)
		panic(conflictErrors)
	}

	_, exist = bidService.bidRepository.FindBidByCorporationAndRequestID(bidService.db, bidInfo.InstallationRequestID, bidInfo.CorporationID)
	if exist {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.AlreadyExist)
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
		Status:           enums.Pending,
	}
	err := bidService.bidRepository.CreateBid(bidService.db, bid)
	if err != nil {
		panic(err)
	}
}

func (bidService *BidService) CancelBid(bidInfo biddto.CancelBidRequest) {
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
	case request.Status != enums.Open:
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
	case bid.Status != enums.Pending:
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.NotExist)
		panic(conflictErrors)
	}
	err := bidService.bidRepository.DeleteBidByID(bidService.db, bidInfo.BidID)
	if err != nil {
		panic(err)
	}
}

func (bidService *BidService) GetBids(corporationID uint, page int, pageSize int, sortBy string, ascending bool) []biddto.BidsResponse {
	offset := (page - 1) * pageSize
	dir := "ASC"
	if !ascending {
		dir = "DESC"
	}

	_, exist := bidService.corporationService.GetCorporationByID(corporationID)
	var conflictErrors exception.ConflictErrors
	if !exist {
		conflictErrors.Add(bidService.constants.Field.Corporation, bidService.constants.Tag.NotRegistered)
		panic(conflictErrors)
	}

	if sortBy == "" {
		sortBy = "created_at"
	}

	bids := bidService.bidRepository.GetBids(bidService.db, corporationID, offset, pageSize, sortBy, dir)
	bidResponses := make([]biddto.BidsResponse, len(bids))
	for i, bid := range bids {
		bidResponses[i] = biddto.BidsResponse{
			ID:                    bid.ID,
			InstallationRequestID: bid.RequestID,
			MinCost:               bid.MinCost,
			MaxCost:               bid.MaxCost,
			MinDeadline:           bid.MinDeadline,
			MaxDeadline:           bid.MaxDeadline,
			Description:           bid.Description,
			InstallationTime:      bid.InstallationTime,
			Status:                bid.Status.String(),
		}
	}

	return bidResponses
}
