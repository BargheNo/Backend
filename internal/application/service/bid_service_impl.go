package serviceimpl

import (
	"encoding/json"
	"log"

	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/message"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
)

type BidService struct {
	constants           *bootstrap.Constants
	installationService service.InstallationService
	userService         service.UserService
	corporationService  service.CorporationService
	paymentService      service.PaymentService
	guaranteeService    service.GuaranteeService
	rabbitMQ            message.Broker
	bidRepository       repository.BidRepository
	db                  database.Database
}

type BidServiceDeps struct {
	Constants           *bootstrap.Constants
	InstallationService service.InstallationService
	UserService         service.UserService
	CorporationService  service.CorporationService
	PaymentService      service.PaymentService
	GuaranteeService    service.GuaranteeService
	RabbitMQ            message.Broker
	BidRepository       repository.BidRepository
	DB                  database.Database
}

func NewBidService(deps BidServiceDeps) *BidService {
	return &BidService{
		constants:           deps.Constants,
		installationService: deps.InstallationService,
		userService:         deps.UserService,
		corporationService:  deps.CorporationService,
		paymentService:      deps.PaymentService,
		guaranteeService:    deps.GuaranteeService,
		rabbitMQ:            deps.RabbitMQ,
		bidRepository:       deps.BidRepository,
		db:                  deps.DB,
	}
}

func (bidService *BidService) GetBidStatuses() []biddto.GetBidStatusesResponse {
	statuses := enum.GetAllBidStatuses()
	response := make([]biddto.GetBidStatusesResponse, len(statuses))
	for i, status := range statuses {
		response[i] = biddto.GetBidStatusesResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response
}

func (bidService *BidService) GetRequestAnonymousBids(requestInfo biddto.GetListRequestBidsRequest) []biddto.AnonymousBidResponse {
	if _, err := bidService.installationService.ValidateRequestOwnership(requestInfo.RequestID, requestInfo.UserID); err != nil {
		panic(err)
	}

	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	allowedStatus := []enum.BidStatus{enum.BidStatusPending, enum.BidStatusAccepted, enum.BidStatusRejected}

	bids := bidService.bidRepository.FindRequestBids(bidService.db, requestInfo.RequestID, allowedStatus, paginationModifier, sortingModifier)
	bidResponses := make([]biddto.AnonymousBidResponse, len(bids))

	for i, bid := range bids {
		paymentTerms, err := bidService.paymentService.GetPaymentTerms(bid.PaymentTermsID)
		if err != nil {
			panic(err)
		}

		var guarantee guaranteedto.GuaranteeResponse
		if bid.GuaranteeID != nil {
			guarantee, err = bidService.guaranteeService.GetGuarantee(*bid.GuaranteeID)
			if err != nil {
				panic(err)
			}
		}

		bidResponses[i] = biddto.AnonymousBidResponse{
			ID:               bid.ID,
			Description:      bid.Description,
			Cost:             bid.Cost,
			InstallationTime: bid.InstallationTime,
			Status:           bid.Status.String(),
			PaymentTerms:     paymentTerms,
			Guarantee:        guarantee,
		}
	}

	return bidResponses
}

func (bidService *BidService) GetRequestBidsByAdmin(requestInfo biddto.GetListRequestBidsRequestByAdmin) []biddto.AdminBidResponse {
	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	allowedStatus := enum.GetAllBidStatuses()

	bids := bidService.bidRepository.FindRequestBids(bidService.db, requestInfo.RequestID, allowedStatus, paginationModifier, sortingModifier)
	bidResponses := make([]biddto.AdminBidResponse, len(bids))

	for i, bid := range bids {
		paymentTerms, err := bidService.paymentService.GetPaymentTerms(bid.PaymentTermsID)
		if err != nil {
			panic(err)
		}

		var guarantee guaranteedto.GuaranteeResponse
		if bid.GuaranteeID != nil {
			guarantee, err = bidService.guaranteeService.GetGuarantee(*bid.GuaranteeID)
			if err != nil {
				panic(err)
			}
		}

		bidder, _ := bidService.userService.GetUserCredential(bid.BidderID)
		corporation := bidService.corporationService.GetCorporationCredentials(bid.CorporationID)

		bidResponses[i] = biddto.AdminBidResponse{
			ID:               bid.ID,
			Corporation:      corporation,
			Bidder:           bidder,
			Description:      bid.Description,
			Cost:             bid.Cost,
			InstallationTime: bid.InstallationTime,
			Status:           bid.Status.String(),
			PaymentTerms:     paymentTerms,
			Guarantee:        guarantee,
		}
	}

	return bidResponses
}

func (bidService *BidService) GetRequestAnonymousBid(requestInfo biddto.GetCustomerBidRequest) biddto.AnonymousBidResponse {
	if _, err := bidService.installationService.ValidateRequestOwnership(requestInfo.RequestID, requestInfo.UserID); err != nil {
		panic(err)
	}

	bid, exist := bidService.bidRepository.FindRequestBid(bidService.db, requestInfo.BidID, requestInfo.RequestID)
	if !exist || bid.Status == enum.BidStatusCanceled {
		notFoundError := exception.NotFoundError{Item: bidService.constants.Field.Bid}
		panic(notFoundError)
	}

	paymentTerms, err := bidService.paymentService.GetPaymentTerms(bid.PaymentTermsID)
	if err != nil {
		panic(err)
	}

	var guarantee guaranteedto.GuaranteeResponse
	if bid.GuaranteeID != nil {
		guarantee, err = bidService.guaranteeService.GetGuarantee(*bid.GuaranteeID)
		if err != nil {
			panic(err)
		}
	}

	return biddto.AnonymousBidResponse{
		ID:               bid.ID,
		Description:      bid.Description,
		Cost:             bid.Cost,
		InstallationTime: bid.InstallationTime,
		Status:           bid.Status.String(),
		PaymentTerms:     paymentTerms,
		Guarantee:        guarantee,
	}
}

// TODO: operator validation will kill us NO NEED TO VALIDATE OPERATOR HERE!!! but ok :)
func (bidService *BidService) AcceptBid(request biddto.GetCustomerBidRequest) {
	installationRequest, err := bidService.installationService.ValidateRequestOwnership(request.RequestID, request.UserID)
	if err != nil {
		panic(err)
	}

	if installationRequest.Status != enum.InstallationRequestStatusActive.String() {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(bidService.constants.Field.InstallationRequest, bidService.constants.Tag.NotActive)
		panic(conflictErrors)
	}

	bid, exist := bidService.bidRepository.FindRequestBid(bidService.db, request.BidID, request.RequestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: bidService.constants.Field.Bid}
		panic(notFoundError)
	}

	if bid.Status != enum.BidStatusPending {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.NotActive)
		panic(conflictErrors)
	}

	changeRequestStatus := installationdto.ChangeRequestStatusRequest{
		OwnerID:   request.UserID,
		Status:    enum.InstallationRequestStatusDone,
		RequestID: request.RequestID,
	}
	bidService.installationService.ChangeInstallationRequestStatus(changeRequestStatus)

	bid.Status = enum.BidStatusAccepted
	if err := bidService.bidRepository.UpdateBid(bidService.db, bid); err != nil {
		panic(err)
	}

	panelInfo := installationdto.AddPanelRequest{
		Name:                 installationRequest.Name,
		Status:               enum.PanelStatusPending,
		CorporationID:        bid.CorporationID,
		OperatorID:           bid.BidderID,
		CustomerPhone:        installationRequest.Customer.Phone,
		Power:                bid.Power,
		Area:                 bid.Area,
		BuildingType:         enum.PanelStatusPending,
		Tilt:                 0,
		Azimuth:              0,
		TotalNumberOfModules: 0,
		Address: addressdto.CreateAddressRequest{
			ProvinceID:    installationRequest.Address.ProvinceID,
			CityID:        installationRequest.Address.CityID,
			StreetAddress: installationRequest.Address.StreetAddress,
			PostalCode:    installationRequest.Address.PostalCode,
			HouseNumber:   installationRequest.Address.HouseNumber,
			Unit:          installationRequest.Address.Unit,
		},
	}
	bidService.installationService.AddPanel(panelInfo)
}

func (bidService *BidService) RejectBid(request biddto.GetCustomerBidRequest) {
	installationRequest, err := bidService.installationService.ValidateRequestOwnership(request.RequestID, request.UserID)
	if err != nil {
		panic(err)
	}

	if installationRequest.Status != enum.InstallationRequestStatusActive.String() {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(bidService.constants.Field.InstallationRequest, bidService.constants.Tag.NotActive)
		panic(conflictErrors)
	}

	bid, exist := bidService.bidRepository.FindRequestBid(bidService.db, request.BidID, request.RequestID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: bidService.constants.Field.Bid}
		panic(notFoundError)
	}

	if bid.Status != enum.BidStatusPending {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.NotActive)
		panic(conflictErrors)
	}

	bid.Status = enum.BidStatusAccepted
	if err := bidService.bidRepository.UpdateBid(bidService.db, bid); err != nil {
		panic(err)
	}
}

func (bidService *BidService) SetBid(bidInfo biddto.SetBidRequest) {
	var conflictErrors exception.ConflictErrors

	if !bidService.corporationService.ISCorporationApproved(bidInfo.CorporationID) {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: bidService.constants.Field.Bid,
		}
		panic(forbiddenError)
	}

	bidService.userService.IsUserActive(bidInfo.BidderID)
	bidService.corporationService.CheckApplicantAccess(bidInfo.CorporationID, bidInfo.BidderID)

	installationRequest := bidService.installationService.GetPublicInstallationRequest(bidInfo.RequestID)
	if installationRequest.Status != enum.InstallationRequestStatusActive.String() {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.ForbiddenStatus)
		panic(conflictErrors)
	}

	allowedStatus := []enum.BidStatus{enum.BidStatusPending}
	_, exist := bidService.bidRepository.FindBidByCorporationAndRequestID(bidService.db, bidInfo.RequestID, bidInfo.CorporationID, allowedStatus)
	if exist {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.AlreadyExist)
		panic(conflictErrors)
	}

	if bidInfo.GuaranteeID != nil {
		if err := bidService.guaranteeService.ValidateActiveGuaranteeOwnerShip(*bidInfo.GuaranteeID, bidInfo.CorporationID); err != nil {
			panic(err)
		}
	}

	paymentTermsID := bidService.paymentService.CreatePaymentTerms(bidInfo.PaymentTerms)

	bid := &entity.Bid{
		CorporationID:    bidInfo.CorporationID,
		BidderID:         bidInfo.BidderID,
		RequestID:        bidInfo.RequestID,
		Status:           bidInfo.Status,
		Cost:             bidInfo.Cost,
		Area:             bidInfo.Area,
		Power:            bidInfo.Power,
		Description:      bidInfo.Description,
		InstallationTime: bidInfo.InstallationTime,
		PaymentTermsID:   paymentTermsID,
		GuaranteeID:      bidInfo.GuaranteeID,
	}
	if err := bidService.bidRepository.CreateBid(bidService.db, bid); err != nil {
		panic(err)
	}

	additionalData := biddto.BidNotificationData{
		BidID: bid.ID,
	}
	data, err := json.Marshal(additionalData)
	if err != nil {
		log.Println("Invalid data for message notification")
	}

	msg := struct {
		TypeName    enum.NotificationType `json:"typeName"`
		RecipientID uint                  `json:"recipientID"`
		Data        []byte                `json:"data"`
	}{
		TypeName:    enum.CorpSendBidNotificationType,
		RecipientID: installationRequest.Customer.ID,
		Data:        data,
	}

	if err := bidService.rabbitMQ.PublishMessage(bidService.constants.RabbitMQ.Events.SendNotification, msg); err != nil {
		log.Printf("error during send notification after bid: %v", err)
	}
}

func (bidService *BidService) GetCorporationBids(request biddto.GetCorporationBidsRequest) []biddto.CorporationBidResponse {
	bidService.corporationService.CheckApplicantAccess(request.CorporationID, request.UserID)

	paginationModifier := repositoryimpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("updated_at", true)

	allowedStatus := []enum.BidStatus{enum.BidStatus(request.Status)}
	if enum.BidStatus(request.Status) == enum.BidStatusAll {
		allowedStatus = enum.GetAllBidStatuses()
	}

	bids := bidService.bidRepository.FindCorporationBids(bidService.db, request.CorporationID, allowedStatus, paginationModifier, sortingModifier)
	bidResponses := make([]biddto.CorporationBidResponse, len(bids))

	for i, bid := range bids {
		request := installationdto.CorporationPanelRequest{
			CorporationID:  request.CorporationID,
			OperatorID:     request.UserID,
			InstallationID: bid.RequestID,
		}
		installationRequest := bidService.installationService.GetAnonymousInstallationRequest(request)
		bidder := bidService.userService.GetUserCredential(bid.BidderID)
		payment, _ := bidService.paymentService.GetPaymentTerms(bid.PaymentTermsID)
		var guarantee guaranteedto.GuaranteeResponse
		if bid.GuaranteeID != nil {
			guarantee, _ = bidService.guaranteeService.GetGuarantee(*bid.GuaranteeID)
		}
		bidResponses[i] = biddto.CorporationBidResponse{
			ID:                  bid.ID,
			Bidder:              bidder,
			InstallationRequest: installationRequest,
			Description:         bid.Description,
			Cost:                bid.Cost,
			InstallationTime:    bid.InstallationTime,
			Status:              bid.Status.String(),
			PaymentTerms:        payment,
			Guarantee:           guarantee,
		}
	}

	return bidResponses
}

func (bidService *BidService) GetCorporationBid(request biddto.GetBidRequest) biddto.CorporationBidResponse {
	bidService.corporationService.CheckApplicantAccess(request.CorporationID, request.UserID)

	bid, exist := bidService.bidRepository.FindCorporationBid(bidService.db, request.BidID, request.CorporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: bidService.constants.Field.Bid}
		panic(notFoundError)
	}

	getInstallationRequest := installationdto.CorporationPanelRequest{
		CorporationID:  request.CorporationID,
		OperatorID:     request.UserID,
		InstallationID: bid.RequestID,
	}
	installationRequest := bidService.installationService.GetAnonymousInstallationRequest(getInstallationRequest)

	bidder := bidService.userService.GetUserCredential(bid.BidderID)
	payment, _ := bidService.paymentService.GetPaymentTerms(bid.PaymentTermsID)

	var guarantee guaranteedto.GuaranteeResponse
	if bid.GuaranteeID != nil {
		guarantee, _ = bidService.guaranteeService.GetGuarantee(*bid.GuaranteeID)
	}

	return biddto.CorporationBidResponse{
		ID:                  bid.ID,
		Bidder:              bidder,
		InstallationRequest: installationRequest,
		Description:         bid.Description,
		Cost:                bid.Cost,
		InstallationTime:    bid.InstallationTime,
		Status:              bid.Status.String(),
		PaymentTerms:        payment,
		Guarantee:           guarantee,
	}
}

func (bidService *BidService) UpdateBid(request biddto.UpdateBidRequest) {
	bidService.corporationService.CheckApplicantAccess(request.CorporationID, request.BidderID)

	bid, exist := bidService.bidRepository.FindCorporationBid(bidService.db, request.BidID, request.CorporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: bidService.constants.Field.Bid}
		panic(notFoundError)
	}

	var conflictErrors exception.ConflictErrors
	if bid.Status == enum.BidStatusAccepted {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.AlreadyAccepted)
		panic(conflictErrors)
	} else if bid.Status == enum.BidStatusCanceled {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.AlreadyCanceled)
		panic(conflictErrors)
	} else if bid.Status != enum.BidStatusPending {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.ForbiddenStatus)
		panic(conflictErrors)
	}

	if request.Cost != nil {
		bid.Cost = *request.Cost
	}

	if request.Area != nil {
		bid.Area = *request.Area
	}

	if request.Power != nil {
		bid.Power = *request.Power
	}

	if request.Description != nil {
		bid.Description = *request.Description
	}

	if request.InstallationTime != nil {
		bid.InstallationTime = *request.InstallationTime
	}

	if request.GuaranteeID != nil {
		bid.GuaranteeID = request.GuaranteeID
	}

	if request.PaymentTerms != nil {
		request.PaymentTerms.ID = bid.PaymentTermsID
		if err := bidService.paymentService.UpdatePaymentTerms(*request.PaymentTerms); err != nil {
			panic(err)
		}
	}

	if err := bidService.bidRepository.UpdateBid(bidService.db, bid); err != nil {
		panic(err)
	}
}

func (bidService *BidService) CancelBid(bidInfo biddto.GetBidRequest) {
	bidService.corporationService.CheckApplicantAccess(bidInfo.CorporationID, bidInfo.UserID)

	bid, exist := bidService.bidRepository.FindCorporationBid(bidService.db, bidInfo.BidID, bidInfo.CorporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: bidService.constants.Field.Bid}
		panic(notFoundError)
	}

	var conflictErrors exception.ConflictErrors
	if bid.Status == enum.BidStatusCanceled {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.AlreadyCanceled)
		panic(conflictErrors)
	} else if bid.Status != enum.BidStatusPending {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.ForbiddenStatus)
		panic(conflictErrors)
	}

	bid.Status = enum.BidStatusCanceled

	if err := bidService.bidRepository.UpdateBid(bidService.db, bid); err != nil {
		panic(err)
	}
}
