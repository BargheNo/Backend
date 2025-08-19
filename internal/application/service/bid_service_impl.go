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
	sortBy := sortby.BidSortBy(requested)
	if _, ok := allowed[sortBy]; ok {
		return sortBy.DBColumn()
	}
	return sortby.NewsSortByCreatedAt.DBColumn()
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

func (bidService *BidService) mapToStatusWithFallback(requested uint, allowedStatuses, defaultValue []enum.BidStatus) []enum.BidStatus {
	for _, status := range allowedStatuses {
		if requested == uint(status) {
			if status == enum.BidStatusAll {
				return enum.GetCorporationBidStatuses()
			}
			if status == enum.BidStatusAllCustomer {
				return enum.GetUserBidStatuses()
			}
			return []enum.BidStatus{status}
		}

	}
	return defaultValue
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
	for col := range columns {
		response[i] = biddto.GetBidEnumResponse{
			ID:   uint(col),
			Name: col.Name(),
		}
		i++
	}
	return response
}

func (bidService *BidService) GetUserBidStatuses() []biddto.GetBidEnumResponse {
	statuses := enum.GetUserBidStatuses()
	response := make([]biddto.GetBidEnumResponse, len(statuses))
	for i, status := range statuses {
		response[i] = biddto.GetBidEnumResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response
}

func (bidService *BidService) GetCorporationBidStatuses() []biddto.GetBidEnumResponse {
	statuses := enum.GetCorporationBidStatuses()
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

	allUserBidStatuses := enum.GetUserBidStatuses()
	allowedStatus := bidService.mapToStatusWithFallback(requestInfo.Status, allUserBidStatuses, allUserBidStatuses)
	bids, err := bidService.bidRepository.FindRequestBids(bidService.db, requestInfo.RequestID, allowedStatus, options)
	if err != nil {
		return nil, 0, err
	}
	bidResponses := make([]biddto.AnonymousBidResponse, len(bids))

	for i, bid := range bids {
		if err := bidService.corporationService.CheckApplicantAccess(bid.CorporationID, bid.BidderID); err != nil {
			continue
		}

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

func (bidService *BidService) getAdminRequestBidsByQuery(requestID uint, allowedStatus []enum.BidStatus, query string, options *postgres.QueryOptions) ([]*entity.Bid, int64, error) {
	if query == "" {
		bids, err := bidService.bidRepository.FindRequestBids(bidService.db, requestID, allowedStatus, options)
		if err != nil {
			return nil, 0, err
		}
		count, err := bidService.bidRepository.CountRequestBids(bidService.db, requestID, allowedStatus)
		if err != nil {
			return nil, 0, err
		}
		return bids, count, nil
	}

	requests, err := bidService.bidRepository.FindRequestBidsByQuery(bidService.db, requestID, allowedStatus, query, options)
	if err != nil {
		return nil, 0, err
	}
	count, err := bidService.bidRepository.CountRequestBidsByQuery(bidService.db, requestID, allowedStatus, query)
	if err != nil {
		return nil, 0, err
	}
	return requests, count, nil
}

func (bidService *BidService) GetRequestBidsByAdmin(requestInfo biddto.GetListRequestBidsRequestByAdmin) ([]biddto.AdminBidResponse, int64, error) {
	options := postgres.NewQueryOptions().
		WithPagination(requestInfo.Limit, requestInfo.Offset).
		WithSorting(bidService.getSortByColumn(requestInfo.SortBy), requestInfo.Asc)

	allStatuses := enum.GetAllBidStatuses()
	allowedStatus := bidService.mapToStatusWithFallback(requestInfo.Status, allStatuses, allStatuses)

	bids, count, err := bidService.getAdminRequestBidsByQuery(requestInfo.RequestID, allowedStatus, requestInfo.Query, options)
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

	return bidResponses, count, nil
}

func (bidService *BidService) getBidsByAdminByQuery(allowedStatus []enum.BidStatus, query string, options *postgres.QueryOptions) ([]*entity.Bid, int64, error) {
	if query == "" {
		bids, err := bidService.bidRepository.FindBidsByStatus(bidService.db, allowedStatus, options)
		if err != nil {
			return nil, 0, err
		}
		count, err := bidService.bidRepository.CountBidsByStatus(bidService.db, allowedStatus)
		if err != nil {
			return nil, 0, err
		}
		return bids, count, nil
	}

	requests, err := bidService.bidRepository.FindBidsByStatusAndQuery(bidService.db, allowedStatus, query, options)
	if err != nil {
		return nil, 0, err
	}
	count, err := bidService.bidRepository.CountBidsByStatusAndQuery(bidService.db, allowedStatus, query)
	if err != nil {
		return nil, 0, err
	}
	return requests, count, nil
}

func (bidService *BidService) GetBidsByAdmin(requestInfo biddto.GetListBidsRequestByAdmin) ([]biddto.AdminBidResponse, int64, error) {
	options := postgres.NewQueryOptions().
		WithPagination(requestInfo.Limit, requestInfo.Offset).
		WithSorting(bidService.getSortByColumn(requestInfo.SortBy), requestInfo.Asc)

	allStatuses := enum.GetAllBidStatuses()
	allowedStatus := bidService.mapToStatusWithFallback(requestInfo.Status, allStatuses, allStatuses)
	bids, count, err := bidService.getBidsByAdminByQuery(allowedStatus, requestInfo.Query, options)
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

	return bidResponses, count, nil
}

func (bidService *BidService) GetBidByAdmin(bidID uint) (biddto.AdminBidResponse, error) {
	bid, err := bidService.bidRepository.FindBidByID(bidService.db, bidID)
	if err != nil {
		return biddto.AdminBidResponse{}, err
	}
	if bid == nil {
		return biddto.AdminBidResponse{}, exception.NotFoundError{Item: bidService.constants.Field.Bid}
	}

	paymentTerms, err := bidService.paymentService.GetPaymentTerms(bid.PaymentTermsID)
	if err != nil {
		return biddto.AdminBidResponse{}, err
	}

	var guarantee guaranteedto.GuaranteeResponse
	if bid.GuaranteeID != nil {
		guarantee, err = bidService.guaranteeService.GetGuarantee(*bid.GuaranteeID)
		if err != nil {
			return biddto.AdminBidResponse{}, err
		}
	}

	bidder, err := bidService.userService.GetUserCredential(bid.BidderID)
	if err != nil {
		return biddto.AdminBidResponse{}, err
	}

	corporation, err := bidService.corporationService.GetCorporationCredentials(bid.CorporationID)
	if err != nil {
		return biddto.AdminBidResponse{}, err
	}

	return biddto.AdminBidResponse{
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
	}, nil
}

func (bidService *BidService) GetRequestAnonymousBid(requestInfo biddto.GetCustomerBidRequest) (biddto.AnonymousBidResponse, error) {
	if _, err := bidService.installationService.ValidateRequestOwnership(requestInfo.RequestID, requestInfo.UserID); err != nil {
		return biddto.AnonymousBidResponse{}, err
	}

	bid, err := bidService.getRequestBid(requestInfo.BidID, requestInfo.RequestID)
	if err != nil {
		return biddto.AnonymousBidResponse{}, err
	}

	if err := bidService.corporationService.CheckApplicantAccess(bid.CorporationID, bid.BidderID); err != nil {
		return biddto.AnonymousBidResponse{}, err
	}

	if bid.Status == enum.BidStatusCanceled || bid.Status == enum.BidStatusExpired {
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

func (bidService *BidService) AcceptBid(request biddto.GetCustomerBidRequest) error {
	if err := bidService.userService.IsUserActive(request.UserID); err != nil {
		return err
	}

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

	bid.Status = enum.BidStatusAccepted

	panelInfo := installationdto.AddPanelRequest{
		Name:          installationRequest.Name,
		Status:        enum.PanelStatusPending,
		CorporationID: bid.CorporationID,
		OperatorID:    bid.BidderID,
		CustomerPhone: installationRequest.Customer.Phone,
		Power:         bid.Power,
		Area:          bid.Area,
		BuildingType:  enum.PanelStatusPending,
		Address: addressdto.CreateAddressRequest{
			ProvinceID:    installationRequest.Address.ProvinceID,
			CityID:        installationRequest.Address.CityID,
			StreetAddress: installationRequest.Address.StreetAddress,
			PostalCode:    installationRequest.Address.PostalCode,
			HouseNumber:   installationRequest.Address.HouseNumber,
			Unit:          installationRequest.Address.Unit,
		},
	}

	err = bidService.db.WithTransaction(func(tx database.Database) error {
		if err := bidService.installationService.ChangeInstallationRequestStatus(changeRequestStatus); err != nil {
			return err
		}

		if err := bidService.bidRepository.UpdateBid(tx, bid); err != nil {
			return err
		}

		if err := bidService.installationService.AddPanel(panelInfo); err != nil {
			return err
		}
		return nil
	})

	return err
}

func (bidService *BidService) RejectBid(request biddto.GetCustomerBidRequest) error {
	var conflictErrors exception.ConflictErrors

	installationRequest, err := bidService.installationService.ValidateRequestOwnership(request.RequestID, request.UserID)
	if err != nil {
		return err
	}

	if installationRequest.Status != enum.InstallationRequestStatusActive.String() {
		conflictErrors.Add(bidService.constants.Field.InstallationRequest, bidService.constants.Tag.NotActive)
		return conflictErrors
	}

	bid, err := bidService.getRequestBid(request.BidID, request.RequestID)
	if err != nil {
		return err
	}

	if bid.Status != enum.BidStatusPending {
		conflictErrors.Add(bidService.constants.Field.Bid, bidService.constants.Tag.NotActive)
		return conflictErrors
	}

	bid.Status = enum.BidStatusRejected
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
	var conflictErrors exception.ConflictErrors
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
		conflictErrors.Add(bidService.constants.Field.InstallationRequest, bidService.constants.Tag.ForbiddenStatus)
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

	if err := bidService.paymentService.ValidatePaymentMethod(bidInfo.PaymentTerms.PaymentMethod); err != nil {
		return err
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

func (bidService *BidService) getCorporationBidsByQuery(corporationID uint, allowedStatus []enum.BidStatus, query string, options *postgres.QueryOptions) ([]*entity.Bid, int64, error) {
	if query == "" {
		bids, err := bidService.bidRepository.FindCorporationBids(bidService.db, corporationID, allowedStatus, options)
		if err != nil {
			return nil, 0, err
		}
		count, err := bidService.bidRepository.CountCorporationBids(bidService.db, corporationID, allowedStatus)
		if err != nil {
			return nil, 0, err
		}
		return bids, count, nil
	}

	requests, err := bidService.bidRepository.FindCorporationBidsByQuery(bidService.db, corporationID, allowedStatus, query, options)
	if err != nil {
		return nil, 0, err
	}
	count, err := bidService.bidRepository.CountCorporationBidsByQuery(bidService.db, corporationID, allowedStatus, query)
	if err != nil {
		return nil, 0, err
	}
	return requests, count, nil
}

func (bidService *BidService) GetCorporationBids(request biddto.GetCorporationBidsRequest) ([]biddto.CorporationBidResponse, int64, error) {
	if err := bidService.corporationService.CheckApplicantAccess(request.CorporationID, request.UserID); err != nil {
		return nil, 0, err
	}

	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset).
		WithSorting(bidService.getSortByColumn(request.SortBy), request.Asc)

	corporationBidStatuses := enum.GetCorporationBidStatuses()
	allowedStatus := bidService.mapToStatusWithFallback(request.Status, corporationBidStatuses, corporationBidStatuses)

	bids, count, err := bidService.getCorporationBidsByQuery(request.CorporationID, allowedStatus, request.Query, options)
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

func (bidService *BidService) applyBidUpdates(bid *entity.Bid, cost, area, power *uint, description *string, installationTime *time.Time) {
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

	bidService.applyBidUpdates(bid, request.Cost, request.Area, request.Power, request.Description, request.InstallationTime)

	if request.GuaranteeID != nil {
		if err := bidService.guaranteeService.ValidateActiveGuaranteeOwnerShip(*request.GuaranteeID, bid.CorporationID); err != nil {
			return err
		}
		bid.GuaranteeID = request.GuaranteeID
	}

	err = bidService.db.WithTransaction(func(tx database.Database) error {
		if request.PaymentTerms != nil {
			if err := bidService.paymentService.ValidatePaymentMethod(*request.PaymentTerms.PaymentMethod); err != nil {
				return err
			}
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
	if err := bidService.corporationService.CheckApplicantAccess(bidInfo.CorporationID, bidInfo.UserID); err != nil {
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

func (bidService *BidService) DeleteBidByAdmin(bidID uint) error {
	bid, err := bidService.bidRepository.FindBidByID(bidService.db, bidID)
	if err != nil {
		return err
	}
	if bid == nil {
		return exception.NotFoundError{Item: bidService.constants.Field.Bid}
	}

	if err := bidService.bidRepository.DeleteBidByID(bidService.db, bidID); err != nil {
		return err
	}
	return nil
}

func (bidService *BidService) UpdateBidByAdmin(request biddto.UpdateBidRequest) error {
	bid, err := bidService.bidRepository.FindBidByID(bidService.db, request.BidID)
	if err != nil {
		return err
	}
	if bid == nil {
		return exception.NotFoundError{Item: bidService.constants.Field.Bid}
	}

	if err := bidService.checkUpdateBidStatus(bid.Status); err != nil {
		return err
	}

	bidService.applyBidUpdates(bid, request.Cost, request.Area, request.Power, request.Description, request.InstallationTime)

	if request.GuaranteeID != nil {
		if err := bidService.guaranteeService.ValidateActiveGuaranteeOwnerShip(*request.GuaranteeID, bid.CorporationID); err != nil {
			return err
		}
		bid.GuaranteeID = request.GuaranteeID
	}

	if err := bidService.bidRepository.UpdateBid(bidService.db, bid); err != nil {
		return err
	}
	return nil
}
