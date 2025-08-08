package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/enum/sortby"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/message"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type BidService struct {
	constants           *bootstrap.Constants
	installationService usecase.InstallationService
	userService         usecase.UserService
	corporationService  usecase.CorporationService
	paymentService      usecase.PaymentService
	guaranteeService    usecase.GuaranteeService
	rabbitMQ            message.Broker
	bidRepository       postgres.BidRepository
	db                  database.Database
}

type BidServiceDeps struct {
	Constants           *bootstrap.Constants
	InstallationService usecase.InstallationService
	UserService         usecase.UserService
	CorporationService  usecase.CorporationService
	PaymentService      usecase.PaymentService
	GuaranteeService    usecase.GuaranteeService
	RabbitMQ            message.Broker
	BidRepository       postgres.BidRepository
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

func (bidService *BidService) getSortByColumn(requested uint) string {
	allowed := sortby.GetBidSortableColumns()
	_, ok := allowed[sortby.BidSortBy(requested)]
	if ok {
		return sortby.BidSortByCreatedAt.DBColumn()
	}
	return sortby.BidSortBy(requested).DBColumn()
}

func (bidService *BidService) getRequestBid(bidID, requestID uint) (*entity.Bid, error) {
	bid, err := bidService.bidRepository.FindRequestBid(bidService.db, bidID, requestID)
	if err != nil {
		return nil, err
	}
	if bid == nil {
		notFoundError := exception.NotFoundError{Item: bidService.constants.Field.Bid}
		return nil, notFoundError
	}
	return bid, nil
}

func (bidService *BidService) getCorporationRequestBid(corporationID, requestID uint, allowedStatus []enum.BidStatus) (*entity.Bid, error) {
	bid, err := bidService.bidRepository.FindBidByCorporationAndRequestID(bidService.db, requestID, corporationID, allowedStatus)
	if err != nil {
		return nil, err
	}
	if bid != nil {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.AlreadyExist)
		return nil, conflictErrors
	}
	return bid, nil
}

func (bidService *BidService) getCorporationBid(bidID, corporationID uint) (*entity.Bid, error) {
	bid, err := bidService.bidRepository.FindCorporationBid(bidService.db, bidID, corporationID)
	if err != nil {
		return nil, err
	}
	if bid == nil {
		notFoundError := exception.NotFoundError{Item: bidService.constants.Field.Bid}
		return nil, notFoundError
	}
	return bid, nil
}

func (bidService *BidService) GetBidSortableColumns() []biddto.GetBidEnumResponse {
	columns := sortby.GetBidSortableColumns()
	response := make([]biddto.GetBidEnumResponse, len(columns))
	i := 0
	for col, _ := range columns {
		response[i] = biddto.GetBidEnumResponse{
			ID:   uint(col),
			Name: col.Name(),
		}
		i++
	}
	return response
}

func (bidService *BidService) GetBidStatuses() []biddto.GetBidEnumResponse {
	statuses := enum.GetAllBidStatuses()
	response := make([]biddto.GetBidEnumResponse, len(statuses))
	for i, status := range statuses {
		response[i] = biddto.GetBidEnumResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response
}

func (bidService *BidService) GetRequestAnonymousBids(requestInfo biddto.GetListRequestBidsRequest) ([]biddto.AnonymousBidResponse, int64, error) {
	if _, err := bidService.installationService.ValidateRequestOwnership(requestInfo.RequestID, requestInfo.UserID); err != nil {
		return nil, 0, err
	}

	options := postgres.NewQueryOptions().
		WithPagination(requestInfo.Limit, requestInfo.Offset).
		WithSorting(bidService.getSortByColumn(requestInfo.SortBy), requestInfo.Asc)

	allowedStatus := []enum.BidStatus{enum.BidStatusPending, enum.BidStatusAccepted, enum.BidStatusRejected}

	bids, err := bidService.bidRepository.FindRequestBids(bidService.db, requestInfo.RequestID, allowedStatus, options)
	if err != nil {
		return nil, 0, err
	}
	bidResponses := make([]biddto.AnonymousBidResponse, len(bids))

	for i, bid := range bids {
		paymentTerms, err := bidService.paymentService.GetPaymentTerms(bid.PaymentTermsID)
		if err != nil {
			return nil, 0, err
		}

		var guarantee guaranteedto.GuaranteeResponse
		if bid.GuaranteeID != nil {
			guarantee, err = bidService.guaranteeService.GetGuarantee(*bid.GuaranteeID)
			if err != nil {
				return nil, 0, err
			}
		}

		bidResponses[i] = biddto.AnonymousBidResponse{
			ID:               bid.ID,
			Description:      bid.Description,
			Area:             bid.Area,
			Power:            bid.Power,
			Cost:             bid.Cost,
			InstallationTime: bid.InstallationTime,
			Status:           bid.Status.String(),
			PaymentTerms:     paymentTerms,
			Guarantee:        guarantee,
		}
	}

	count, err := bidService.bidRepository.CountRequestBids(bidService.db, requestInfo.RequestID, allowedStatus)
	if err != nil {
		return nil, 0, err
	}

	return bidResponses, count, nil
}

func (bidService *BidService) GetRequestBidsByAdmin(requestInfo biddto.GetListRequestBidsRequestByAdmin) ([]biddto.AdminBidResponse, int64, error) {
	allowedStatus := enum.GetAllBidStatuses()

	options := postgres.NewQueryOptions().
		WithPagination(requestInfo.Limit, requestInfo.Offset).
		WithSorting(bidService.getSortByColumn(requestInfo.SortBy), requestInfo.Asc)

	bids, err := bidService.bidRepository.FindRequestBids(bidService.db, requestInfo.RequestID, allowedStatus, options)
	if err != nil {
		return nil, 0, err
	}
	bidResponses := make([]biddto.AdminBidResponse, len(bids))

	for i, bid := range bids {
		paymentTerms, err := bidService.paymentService.GetPaymentTerms(bid.PaymentTermsID)
		if err != nil {
			return nil, 0, err
		}

		var guarantee guaranteedto.GuaranteeResponse
		if bid.GuaranteeID != nil {
			guarantee, err = bidService.guaranteeService.GetGuarantee(*bid.GuaranteeID)
			if err != nil {
				return nil, 0, err
			}
		}

		bidder, err := bidService.userService.GetUserCredential(bid.BidderID)
		if err != nil {
			return nil, 0, err
		}
		corporation, err := bidService.corporationService.GetCorporationCredentials(bid.CorporationID)
		if err != nil {
			return nil, 0, err
		}

		bidResponses[i] = biddto.AdminBidResponse{
			ID:               bid.ID,
			Corporation:      corporation,
			Bidder:           bidder,
			Description:      bid.Description,
			Cost:             bid.Cost,
			Area:             bid.Area,
			Power:            bid.Power,
			InstallationTime: bid.InstallationTime,
			Status:           bid.Status.String(),
			PaymentTerms:     paymentTerms,
			Guarantee:        guarantee,
		}
	}

	count, err := bidService.bidRepository.CountRequestBids(bidService.db, requestInfo.RequestID, allowedStatus)
	if err != nil {
		return nil, 0, err
	}

	return bidResponses, count, nil
}

func (bidService *BidService) GetRequestAnonymousBid(requestInfo biddto.GetCustomerBidRequest) (biddto.AnonymousBidResponse, error) {
	if _, err := bidService.installationService.ValidateRequestOwnership(requestInfo.RequestID, requestInfo.UserID); err != nil {
		return biddto.AnonymousBidResponse{}, err
	}

	bid, err := bidService.getRequestBid(requestInfo.BidID, requestInfo.RequestID)
	if err != nil {
		return biddto.AnonymousBidResponse{}, err
	}

	if bid.Status == enum.BidStatusCanceled {
		notFoundError := exception.NotFoundError{Item: bidService.constants.Field.Bid}
		return biddto.AnonymousBidResponse{}, notFoundError
	}

	paymentTerms, err := bidService.paymentService.GetPaymentTerms(bid.PaymentTermsID)
	if err != nil {
		return biddto.AnonymousBidResponse{}, err
	}

	var guarantee guaranteedto.GuaranteeResponse
	if bid.GuaranteeID != nil {
		guarantee, err = bidService.guaranteeService.GetGuarantee(*bid.GuaranteeID)
		if err != nil {
			return biddto.AnonymousBidResponse{}, err
		}
	}

	return biddto.AnonymousBidResponse{
		ID:               bid.ID,
		Description:      bid.Description,
		Cost:             bid.Cost,
		Area:             bid.Area,
		Power:            bid.Power,
		InstallationTime: bid.InstallationTime,
		Status:           bid.Status.String(),
		PaymentTerms:     paymentTerms,
		Guarantee:        guarantee,
	}, nil
}

// TODO: operator validation will kill us NO NEED TO VALIDATE OPERATOR HERE!!! but ok :)
func (bidService *BidService) AcceptBid(request biddto.GetCustomerBidRequest) error {
	installationRequest, err := bidService.installationService.ValidateRequestOwnership(request.RequestID, request.UserID)
	if err != nil {
		return err
	}

	if installationRequest.Status != enum.InstallationRequestStatusActive.String() {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(bidService.constants.Field.InstallationRequest, bidService.constants.Tag.NotActive)
		return conflictErrors
	}

	bid, err := bidService.getRequestBid(request.BidID, request.RequestID)
	if err != nil {
		return err
	}

	if bid.Status != enum.BidStatusPending {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.NotActive)
		return conflictErrors
	}

	changeRequestStatus := installationdto.ChangeRequestStatusRequest{
		OwnerID:   request.UserID,
		Status:    enum.InstallationRequestStatusDone,
		RequestID: request.RequestID,
	}

	err = bidService.db.WithTransaction(func(tx database.Database) error {
		if err := bidService.installationService.ChangeInstallationRequestStatus(changeRequestStatus); err != nil {
			return err
		}

		bid.Status = enum.BidStatusAccepted
		if err := bidService.bidRepository.UpdateBid(tx, bid); err != nil {
			return err
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
		if err := bidService.installationService.AddPanel(panelInfo); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (bidService *BidService) RejectBid(request biddto.GetCustomerBidRequest) error {
	installationRequest, err := bidService.installationService.ValidateRequestOwnership(request.RequestID, request.UserID)
	if err != nil {
		return err
	}

	if installationRequest.Status != enum.InstallationRequestStatusActive.String() {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(bidService.constants.Field.InstallationRequest, bidService.constants.Tag.NotActive)
		return conflictErrors
	}

	bid, err := bidService.getRequestBid(request.BidID, request.RequestID)
	if err != nil {
		return err
	}

	if bid.Status != enum.BidStatusPending {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.NotActive)
		return conflictErrors
	}

	bid.Status = enum.BidStatusAccepted
	if err := bidService.bidRepository.UpdateBid(bidService.db, bid); err != nil {
		return err
	}
	return nil
}

func (bidService *BidService) sendNotification(requestID, bidID, customerID uint) {
	additionalData := biddto.BidNotificationData{
		RequestID: requestID,
		BidID:     bidID,
		UserID:    customerID,
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
		RecipientID: customerID,
		Data:        data,
	}

	if err := bidService.rabbitMQ.PublishMessage(bidService.constants.RabbitMQ.Events.SendNotification, msg); err != nil {
		log.Printf("error during send notification after bid: %v", err)
	}
}

func (bidService *BidService) SetBid(bidInfo biddto.SetBidRequest) error {
	if err := bidService.userService.IsUserActive(bidInfo.BidderID); err != nil {
		return err
	}

	if err := bidService.corporationService.ISCorporationApproved(bidInfo.CorporationID); err != nil {
		return err
	}

	if err := bidService.corporationService.CheckApplicantAccess(bidInfo.CorporationID, bidInfo.BidderID); err != nil {
		return err
	}

	installationRequest, err := bidService.installationService.GetPublicInstallationRequest(bidInfo.RequestID)
	if err != nil {
		return err
	}

	if installationRequest.Status != enum.InstallationRequestStatusActive.String() {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.ForbiddenStatus)
		return conflictErrors
	}

	allowedStatus := []enum.BidStatus{enum.BidStatusPending}
	bid, err := bidService.getCorporationRequestBid(bidInfo.CorporationID, bidInfo.RequestID, allowedStatus)
	if err != nil {
		return err
	}

	if bidInfo.GuaranteeID != nil {
		if err := bidService.guaranteeService.ValidateActiveGuaranteeOwnerShip(*bidInfo.GuaranteeID, bidInfo.CorporationID); err != nil {
			return err
		}
	}

	err = bidService.db.WithTransaction(func(tx database.Database) error {
		paymentTermsID, err := bidService.paymentService.CreatePaymentTerms(bidInfo.PaymentTerms)
		if err != nil {
			return err
		}

		bid = &entity.Bid{
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
		if err := bidService.bidRepository.CreateBid(tx, bid); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	bidService.sendNotification(bid.RequestID, bid.ID, installationRequest.Customer.ID)

	return nil
}

func (bidService *BidService) GetCorporationBids(request biddto.GetCorporationBidsRequest) ([]biddto.CorporationBidResponse, int64, error) {
	if err := bidService.corporationService.CheckApplicantAccess(request.CorporationID, request.UserID); err != nil {
		return nil, 0, err
	}

	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset).
		WithSorting(bidService.getSortByColumn(request.SortBy), request.Asc)

	allowedStatus := []enum.BidStatus{enum.BidStatus(request.Status)}
	if enum.BidStatus(request.Status) == enum.BidStatusAll {
		allowedStatus = enum.GetAllBidStatuses()
	}

	bids, err := bidService.bidRepository.FindCorporationBids(bidService.db, request.CorporationID, allowedStatus, options)
	if err != nil {
		return nil, 0, err
	}
	bidResponses := make([]biddto.CorporationBidResponse, len(bids))

	for i, bid := range bids {
		request := installationdto.CorporationPanelRequest{
			CorporationID:  request.CorporationID,
			OperatorID:     request.UserID,
			InstallationID: bid.RequestID,
		}
		installationRequest, err := bidService.installationService.GetAnonymousInstallationRequest(request)
		if err != nil {
			return nil, 0, err
		}

		bidder, err := bidService.userService.GetUserCredential(bid.BidderID)
		if err != nil {
			return nil, 0, err
		}

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
			Area:                bid.Area,
			Power:               bid.Power,
			InstallationTime:    bid.InstallationTime,
			Status:              bid.Status.String(),
			PaymentTerms:        payment,
			Guarantee:           guarantee,
		}
	}

	count, err := bidService.bidRepository.CountCorporationBids(bidService.db, request.CorporationID, allowedStatus)
	if err != nil {
		return nil, 0, err
	}

	return bidResponses, count, nil
}

func (bidService *BidService) GetCorporationBid(request biddto.GetBidRequest) (biddto.CorporationBidResponse, error) {
	if err := bidService.corporationService.CheckApplicantAccess(request.CorporationID, request.UserID); err != nil {
		return biddto.CorporationBidResponse{}, err
	}

	bid, err := bidService.getCorporationBid(request.BidID, request.CorporationID)
	if err != nil {
		return biddto.CorporationBidResponse{}, err
	}

	getInstallationRequest := installationdto.CorporationPanelRequest{
		CorporationID:  request.CorporationID,
		OperatorID:     request.UserID,
		InstallationID: bid.RequestID,
	}
	installationRequest, err := bidService.installationService.GetAnonymousInstallationRequest(getInstallationRequest)
	if err != nil {
		return biddto.CorporationBidResponse{}, err
	}

	bidder, err := bidService.userService.GetUserCredential(bid.BidderID)
	if err != nil {
		return biddto.CorporationBidResponse{}, err
	}

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
		Area:                bid.Area,
		Power:               bid.Power,
		InstallationTime:    bid.InstallationTime,
		Status:              bid.Status.String(),
		PaymentTerms:        payment,
		Guarantee:           guarantee,
	}, nil
}

func (bidService *BidService) checkUpdateBidStatus(status enum.BidStatus) error {
	var conflictErrors exception.ConflictErrors
	if status == enum.BidStatusAccepted {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.AlreadyAccepted)
		return conflictErrors
	} else if status == enum.BidStatusCanceled {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.AlreadyCanceled)
		return conflictErrors
	} else if status != enum.BidStatusPending {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.ForbiddenStatus)
		return conflictErrors
	}
	return nil
}

func (bidService *BidService) applyBidUpdates(bid *entity.Bid, cost, area, power, guaranteeID *uint, description *string, installationTime *time.Time) {
	if cost != nil {
		bid.Cost = *cost
	}

	if area != nil {
		bid.Area = *area
	}

	if power != nil {
		bid.Power = *power
	}

	if description != nil {
		bid.Description = *description
	}

	if installationTime != nil {
		bid.InstallationTime = *installationTime
	}

	if guaranteeID != nil {
		bid.GuaranteeID = guaranteeID
	}
}

func (bidService *BidService) UpdateBid(request biddto.UpdateBidRequest) error {
	err := bidService.corporationService.CheckApplicantAccess(request.CorporationID, request.BidderID)
	if err != nil {
		return err
	}

	if err := bidService.userService.IsUserActive(request.BidderID); err != nil {
		return err
	}

	bid, err := bidService.getCorporationBid(request.BidID, request.CorporationID)
	if err != nil {
		return err
	}

	if err := bidService.checkUpdateBidStatus(bid.Status); err != nil {
		return err
	}

	bidService.applyBidUpdates(bid, request.Cost, request.Area, request.Power, request.GuaranteeID, request.Description, request.InstallationTime)

	err = bidService.db.WithTransaction(func(tx database.Database) error {
		if request.PaymentTerms != nil {
			request.PaymentTerms.ID = bid.PaymentTermsID
			if err := bidService.paymentService.UpdatePaymentTerms(*request.PaymentTerms); err != nil {
				return err
			}
		}

		if err := bidService.bidRepository.UpdateBid(tx, bid); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (bidService *BidService) CancelBid(bidInfo biddto.GetBidRequest) error {
	err := bidService.corporationService.CheckApplicantAccess(bidInfo.CorporationID, bidInfo.UserID)
	if err != nil {
		return err
	}

	if err := bidService.userService.IsUserActive(bidInfo.UserID); err != nil {
		return err
	}

	bid, err := bidService.getCorporationBid(bidInfo.BidID, bidInfo.CorporationID)
	if err != nil {
		return err
	}

	var conflictErrors exception.ConflictErrors
	if bid.Status == enum.BidStatusCanceled {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.AlreadyCanceled)
		return conflictErrors
	} else if bid.Status != enum.BidStatusPending {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.ForbiddenStatus)
		return conflictErrors
	}

	bid.Status = enum.BidStatusCanceled

	if err := bidService.bidRepository.UpdateBid(bidService.db, bid); err != nil {
		return err
	}
	return nil
}
