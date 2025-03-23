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
	constants          *bootstrap.Constants
	JWTService         service.JWTService
	db                 database.Database
	bidRepository      repository.BidRepository
	addressService     service.AddressService
	corporationService service.CorporationService
}

func NewBidService(
	constants *bootstrap.Constants,
	jwtService service.JWTService,
	db database.Database,
	bidRepository repository.BidRepository,
	addressService service.AddressService,
	corporationService service.CorporationService,
) *BidService {
	return &BidService{
		constants:          constants,
		JWTService:         jwtService,
		db:                 db,
		bidRepository:      bidRepository,
		addressService:     addressService,
		corporationService: corporationService,
	}
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
	case installationRequest.Status != enum.Active:
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
		Cost:             bidInfo.Cost,
		Description:      bidInfo.Description,
		InstallationDate: bidInfo.InstallationDate,
		Status:           enum.Pending,
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
	case request.Status != enum.Active:
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
	case bid.Status != enum.Pending:
		notFoundError = exception.NotFoundError{Item: bidService.constants.Field.Bid}
		panic(notFoundError)
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
	installationRequests := make([]biddto.InstallationRequestDetails, len(bids))
	for i, bid := range bids {
		installationRequests[i] = biddto.InstallationRequestDetails{
			ID:           bid.Request.ID,
			Name:         bid.Request.Name,
			CustomerName: bid.Request.Owner.FirstName + " " + bid.Request.Owner.LastName,
			Address:      bidService.addressService.GetAddress(bid.Request.AddressID),
			PowerRequest: bid.Request.PowerRequest,
		}
	}
	bidResponses := make([]biddto.BidsResponse, len(bids))
	for i, bid := range bids {
		bidResponses[i] = biddto.BidsResponse{
			ID:                         bid.ID,
			InstallationRequestDetails: installationRequests[i],
			Description:                bid.Description,
			Cost:                       bid.Cost,
			InstallationDate:           bid.InstallationDate,
			Status:                     bid.Status.String(),
		}
	}
	return bidResponses
}
