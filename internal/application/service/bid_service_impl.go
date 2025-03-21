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

func (bidService *BidService) GetInstallationRequests(corporationId uint, page int, pageSize int, sortBy string, dir string) []biddto.InstallationRequestResponse {
	offset := (page - 1) * pageSize
	var order string
	if sortBy == "" {
		order = "RANDOM()"
	} else {
		order = sortBy + " " + dir
	}
	_, exist := bidService.corporationService.GetCorporationByID(corporationId)
	if !exist {
		notFoundError := exception.NotFoundError{Item: bidService.constants.Field.Corporation}
		panic(notFoundError)
	}
	var installationRequests []*entity.InstallationRequest
	installationRequests = bidService.bidRepository.GetInstallationRequests(bidService.db, enums.Open, corporationId, offset, pageSize, order)
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
	var notFoundError exception.NotFoundError
	var conflictErrors exception.ConflictErrors
	if !exist {
		notFoundError = exception.NotFoundError{Item: bidService.constants.Field.Corporation}
		panic(notFoundError)
	}
	installationRequest, exist := bidService.bidRepository.FindInstallationRequestByID(bidService.db, bidInfo.InstallationRequestID)
	switch {
	case !exist:
		notFoundError = exception.NotFoundError{Item: bidService.constants.Field.InstallationRequest}
		panic(notFoundError)
	case installationRequest.Status != enums.Open:
		notFoundError = exception.NotFoundError{Item: bidService.constants.Field.InstallationRequest}
		panic(notFoundError)
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
	var notFoundError exception.NotFoundError
	if !exist {
		notFoundError = exception.NotFoundError{Item: bidService.constants.Field.Corporation}
		panic(notFoundError)
	}

	request, exist := bidService.bidRepository.FindInstallationRequestByID(bidService.db, bidInfo.InstallationRequestID)
	switch {
	case !exist:
		notFoundError = exception.NotFoundError{Item: bidService.constants.Field.InstallationRequest}
		panic(notFoundError)
	case request.Status != enums.Open:
		notFoundError = exception.NotFoundError{Item: bidService.constants.Field.InstallationRequest}
		panic(notFoundError)
	}

	bid, exist := bidService.bidRepository.FindBidByID(bidService.db, bidInfo.BidID)
	switch {
	case !exist:
		notFoundError = exception.NotFoundError{Item: bidService.constants.Field.Bid}
		panic(notFoundError)
	case bid.CorporationID != bidInfo.CorporationID:
		notFoundError = exception.NotFoundError{Item: bidService.constants.Field.Bid}
		panic(notFoundError)
	case bid.Status != enums.Pending:
		notFoundError = exception.NotFoundError{Item: bidService.constants.Field.Bid}
		panic(notFoundError)
	}
	err := bidService.bidRepository.DeleteBidByID(bidService.db, bidInfo.BidID)
	if err != nil {
		panic(err)
	}
}

func (bidService *BidService) GetBids(corporationID uint, page int, pageSize int, sortBy string, dir string) []biddto.BidsResponse {
	offset := (page - 1) * pageSize
	var order string
	if sortBy == "" {
		order = "created_at" + " " + dir
	} else {
		order = sortBy + " " + dir
	}
	_, exist := bidService.corporationService.GetCorporationByID(corporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: bidService.constants.Field.Corporation}
		panic(notFoundError)
	}

	bids := bidService.bidRepository.GetBids(bidService.db, corporationID, offset, pageSize, order)
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
