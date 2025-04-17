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
	userService         service.UserService
	corporationService  service.CorporationService
	bidRepository       repository.BidRepository
	db                  database.Database
}

func NewBidService(
	constants *bootstrap.Constants,
	installationService service.InstallationService,
	userService service.UserService,
	corporationService service.CorporationService,
	bidRepository repository.BidRepository,
	db database.Database,
) *BidService {
	return &BidService{
		constants:           constants,
		installationService: installationService,
		userService:         userService,
		corporationService:  corporationService,
		bidRepository:       bidRepository,
		db:                  db,
	}
}

func (bidService *BidService) SetBid(bidInfo biddto.SetBidRequest) {
	var conflictErrors exception.ConflictErrors
	approved := bidService.corporationService.ISCorporationApproved(bidInfo.CorporationID)
	if !approved {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: bidService.constants.Field.Bid,
		}
		panic(forbiddenError)
	}
	bidService.corporationService.CheckApplicantAccess(bidInfo.CorporationID, bidInfo.BidderID)

	installationRequest := bidService.installationService.GetInstallationRequest(bidInfo.InstallationRequestID)
	if installationRequest.Status != enum.InstallationRequestStatusActive.String() {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.ForbiddenStatus)
		panic(conflictErrors)
	}

	allowedStatus := []enum.BidStatus{enum.BidStatusRejected, enum.BidStatusExpired}
	_, exist := bidService.bidRepository.FindBidByCorporationAndRequestID(bidService.db, bidInfo.InstallationRequestID, bidInfo.CorporationID, allowedStatus)
	if exist {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.AlreadyExist)
		panic(conflictErrors)
	}

	bid := &entity.Bid{
		CorporationID:    bidInfo.CorporationID,
		BidderID:         bidInfo.BidderID,
		RequestID:        bidInfo.InstallationRequestID,
		Cost:             bidInfo.Cost,
		Description:      bidInfo.Description,
		InstallationTime: bidInfo.InstallationTime,
		Status:           enum.BidStatusPending,
	}
	err := bidService.bidRepository.CreateBid(bidService.db, bid)
	if err != nil {
		panic(err)
	}
}

func (bidService *BidService) CancelBid(bidInfo biddto.CancelBidRequest) {
	var conflictErrors exception.ConflictErrors
	approved := bidService.corporationService.ISCorporationApproved(bidInfo.CorporationID)
	if !approved {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: bidService.constants.Field.Bid,
		}
		panic(forbiddenError)
	}

	bidService.corporationService.CheckApplicantAccess(bidInfo.CorporationID, bidInfo.BidderID)

	request := bidService.installationService.GetInstallationRequest(bidInfo.InstallationRequestID)
	if request.Status != enum.InstallationRequestStatusActive.String() {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.ForbiddenStatus)
		panic(conflictErrors)
	}

	bid, exist := bidService.bidRepository.FindBidByID(bidService.db, bidInfo.BidID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: bidService.constants.Field.Bid}
		panic(notFoundError)
	}
	if bid.CorporationID != bidInfo.CorporationID {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: bidService.constants.Field.Bid,
		}
		panic(forbiddenError)
	}
	if bid.Status != enum.BidStatusPending {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.ForbiddenStatus)
		panic(conflictErrors)
	}
	err := bidService.bidRepository.DeleteBidByID(bidService.db, bidInfo.BidID)
	if err != nil {
		panic(err)
	}
}

func (bidService *BidService) GetCorporationBids(bidsRequest biddto.GetCorporationBidsRequest) []biddto.BidsResponse {
	approved := bidService.corporationService.ISCorporationApproved(bidsRequest.CorporationID)
	if !approved {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: bidService.constants.Field.Bid,
		}
		panic(forbiddenError)
	}

	bidService.corporationService.CheckApplicantAccess(bidsRequest.CorporationID, bidsRequest.UserID)

	bids := bidService.bidRepository.FindCorporationBids(bidService.db, bidsRequest.CorporationID, bidsRequest.Offset, bidsRequest.Limit)

	bidResponses := make([]biddto.BidsResponse, len(bids))
	for i, bid := range bids {
		installationRequest := bidService.installationService.GetInstallationRequest(bid.RequestID)
		bidder := bidService.userService.GetUserCredential(bidsRequest.UserID)
		bidResponses[i] = biddto.BidsResponse{
			ID:                         bid.ID,
			Bidder:                     bidder,
			InstallationRequestDetails: installationRequest,
			Description:                bid.Description,
			Cost:                       bid.Cost,
			InstallationTime:           bid.InstallationTime,
			Status:                     bid.Status.String(),
		}
	}

	return bidResponses
}

func (bidService *BidService) GetRequestBids(requestInfo biddto.GetRequestBidsRequest) []biddto.BidsResponse {
	request := bidService.installationService.GetInstallationRequestModel(requestInfo.RequestID)
	if request.OwnerID != requestInfo.UserID {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: bidService.constants.Field.InstallationRequest,
		}
		panic(forbiddenError)
	}

	bids := bidService.bidRepository.FindRequestBids(bidService.db, requestInfo.RequestID)
	bidResponses := make([]biddto.BidsResponse, len(bids))
	for i, bid := range bids {
		installationRequest := bidService.installationService.GetInstallationRequest(bid.RequestID)
		bidder := bidService.userService.GetUserCredential(bid.BidderID)
		bidResponses[i] = biddto.BidsResponse{
			ID:                         bid.ID,
			Bidder:                     bidder,
			InstallationRequestDetails: installationRequest,
			Description:                bid.Description,
			Cost:                       bid.Cost,
			InstallationTime:           bid.InstallationTime,
			Status:                     bid.Status.String(),
		}
	}

	return bidResponses
}
