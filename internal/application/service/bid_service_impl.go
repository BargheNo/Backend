package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type BidService struct {
	constants           *bootstrap.Constants
	installationService service.InstallationService
	jwtService          service.JWTService
	corporationService  service.CorporationService
	bidRepository       repository.BidRepository
	db                  database.Database
}

func NewBidService(
	constants *bootstrap.Constants,
	installationService service.InstallationService,
	jwtService service.JWTService,
	corporationService service.CorporationService,
	bidRepository repository.BidRepository,
	db database.Database,
) *BidService {
	return &BidService{
		constants:           constants,
		installationService: installationService,
		jwtService:          jwtService,
		corporationService:  corporationService,
		bidRepository:       bidRepository,
		db:                  db,
	}
}

func (bidService *BidService) SetBid(bidInfo biddto.SetBidRequest) {
	var notFoundError exception.NotFoundError
	var conflictErrors exception.ConflictErrors

	_, exist := bidService.corporationService.GetCorporationByID(bidInfo.CorporationID)
	if !exist {
		notFoundError = exception.NotFoundError{Item: bidService.constants.Field.Corporation}
		panic(notFoundError)
	}

	installationRequest := bidService.installationService.GetInstallationRequest(bidInfo.InstallationRequestID)
	switch {
	case !exist:
		notFoundError = exception.NotFoundError{Item: bidService.constants.Field.InstallationRequest}
		panic(notFoundError)
	case installationRequest.Status != enum.InstallationRequestStatusActive.String():
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.ForbiddenStatus)
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
		Cost:             bidInfo.Cost,
		Description:      bidInfo.Description,
		InstallationTime: bidInfo.InstallationDate,
		Status:           enum.BidStatusPending,
	}
	err := bidService.bidRepository.CreateBid(bidService.db, bid)
	if err != nil {
		panic(err)
	}
}

func (bidService *BidService) CancelBid(bidInfo biddto.CancelBidRequest) {
	var conflictErrors exception.ConflictErrors

	_, exist := bidService.corporationService.GetCorporationByID(bidInfo.CorporationID)
	var notFoundError exception.NotFoundError
	if !exist {
		notFoundError = exception.NotFoundError{Item: bidService.constants.Field.Corporation}
		panic(notFoundError)
	}
	request := bidService.installationService.GetInstallationRequest(bidInfo.InstallationRequestID)
	if request.Status != enum.InstallationRequestStatusActive.String() {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.ForbiddenStatus)
		panic(conflictErrors)
	}

	bid, exist := bidService.bidRepository.FindBidByID(bidService.db, bidInfo.BidID)
	switch {
	case !exist:
		notFoundError = exception.NotFoundError{Item: bidService.constants.Field.Bid}
		panic(notFoundError)
	case bid.CorporationID != bidInfo.CorporationID:
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: bidService.constants.Field.Bid,
		}
		panic(forbiddenError)
	case bid.Status != enum.BidStatusPending:
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.ForbiddenStatus)
		panic(conflictErrors)
	}
	err := bidService.bidRepository.DeleteBidByID(bidService.db, bidInfo.BidID)
	if err != nil {
		panic(err)
	}
}

func (bidService *BidService) GetBids(bidsRequest biddto.GetBidsRequest) []biddto.BidsResponse {
	_, exist := bidService.corporationService.GetCorporationByID(bidsRequest.CorporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: bidService.constants.Field.Corporation}
		panic(notFoundError)
	}

	bids := bidService.bidRepository.GetBids(bidService.db, bidsRequest.CorporationID, bidsRequest.Offset, bidsRequest.Limit)

	bidResponses := make([]biddto.BidsResponse, len(bids))
	for i, bid := range bids {
		installationRequest := bidService.installationService.GetInstallationRequest(bid.RequestID)
		bidResponses[i] = biddto.BidsResponse{
			ID:                         bid.ID,
			InstallationRequestDetails: installationRequest,
			Description:                bid.Description,
			Cost:                       bid.Cost,
			InstallationTime:           bid.InstallationTime,
			Status:                     bid.Status.String(),
		}
	}

	return bidResponses
}
