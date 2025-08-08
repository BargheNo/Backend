package service

import (
	"encoding/json"

	"github.com/BargheNo/Backend/bootstrap"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	notificationdto "github.com/BargheNo/Backend/internal/application/dto/notification"
	reportdto "github.com/BargheNo/Backend/internal/application/dto/report"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/communication"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/enum/sortby"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/message"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/BargheNo/Backend/internal/infrastructure/websocket"
)

type NotificationService struct {
	constants              *bootstrap.Constants
	userService            usecase.UserService
	bidService             usecase.BidService
	reportService          usecase.ReportService
	emailService           communication.EmailService
	notificationRepository postgres.NotificationRepository
	wsHub                  *websocket.Hub
	rabbitMQ               message.Broker
	db                     database.Database
}

type NotificationServiceDeps struct {
	Constants              *bootstrap.Constants
	UserService            usecase.UserService
	BidService             usecase.BidService
	ReportService          usecase.ReportService
	EmailService           communication.EmailService
	NotificationRepository postgres.NotificationRepository
	WSHub                  *websocket.Hub
	RabbitMQ               message.Broker
	DB                     database.Database
}

func NewNotificationService(deps NotificationServiceDeps) *NotificationService {
	return &NotificationService{
		constants:              deps.Constants,
		userService:            deps.UserService,
		bidService:             deps.BidService,
		reportService:          deps.ReportService,
		emailService:           deps.EmailService,
		notificationRepository: deps.NotificationRepository,
		wsHub:                  deps.WSHub,
		rabbitMQ:               deps.RabbitMQ,
		db:                     deps.DB,
	}
}

func (notificationService *NotificationService) getSortByColumn(requested uint) string {
	allowed := sortby.GetNotificationSortableColumns()
	_, ok := allowed[sortby.NotificationSortBy(requested)]
	if ok {
		return sortby.NotificationSortByCreatedAt.DBColumn()
	}
	return sortby.NotificationSortBy(requested).DBColumn()
}

func (notificationService *NotificationService) GetNotificationSortableColumns() []notificationdto.NotificationEnumResponse {
	columns := sortby.GetNotificationSortableColumns()
	response := make([]notificationdto.NotificationEnumResponse, len(columns))
	i := 0
	for col, _ := range columns {
		response[i] = notificationdto.NotificationEnumResponse{
			ID:   uint(col),
			Name: col.Name(),
		}
		i++
	}
	return response
}

func (notificationService *NotificationService) CreateAndSendNotification(typeName enum.NotificationType, recipientID uint, data []byte) error {
	notificationType, err := notificationService.notificationRepository.GetNotificationTypeByName(notificationService.db, typeName)
	if err != nil {
		return err
	}
	if notificationType == nil {
		notFoundError := exception.NotFoundError{Item: notificationService.constants.Field.NotificationType}
		return notFoundError
	}
	notification := &entity.Notification{
		TypeID:      notificationType.ID,
		RecipientID: recipientID,
		Data:        data,
		IsRead:      false,
	}

	if err := notificationService.notificationRepository.CreateNotification(notificationService.db, notification); err != nil {
		return err
	}

	if err := notificationService.SendNotification(notification, notificationType); err != nil {
		return err
	}

	return nil
}

func (notificationService *NotificationService) enrichBidData(rawData []byte) (map[string]interface{}, error) {
	var bidData biddto.BidNotificationData
	var result map[string]interface{}
	if err := json.Unmarshal(rawData, &bidData); err != nil {
		return nil, err
	}

	requestInfo := biddto.GetCustomerBidRequest(bidData)
	bid, err := notificationService.bidService.GetRequestAnonymousBid(requestInfo)
	if err != nil {
		return nil, err
	}

	bidBytes, err := json.Marshal(bid)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bidBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (notificationService *NotificationService) enrichMaintenanceReportData(rawData []byte) (map[string]interface{}, error) {
	var reportData reportdto.ReportNotificationData
	var result map[string]interface{}
	if err := json.Unmarshal(rawData, &reportData); err != nil {
		return nil, err
	}
	report, err := notificationService.reportService.GetMaintenanceReport(reportData.ReportID)
	if err != nil {
		return nil, err
	}

	reportBytes, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(reportBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (notificationService *NotificationService) enrichPanelReportData(rawData []byte) (map[string]interface{}, error) {
	var reportData reportdto.ReportNotificationData
	var result map[string]interface{}
	if err := json.Unmarshal(rawData, &reportData); err != nil {
		return nil, err
	}
	report, err := notificationService.reportService.GetPanelReport(reportData.ReportID)
	if err != nil {
		return nil, err
	}

	reportBytes, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(reportBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (notificationService *NotificationService) dataCatcher(notificationType enum.NotificationType, notificationData []byte) (map[string]interface{}, error) {
	var data map[string]interface{}
	var err error

	switch notificationType {
	case enum.CorpSendBidNotificationType:
		data, err = notificationService.enrichBidData(notificationData)
		if err != nil {
			return nil, err
		}
	case enum.MaintenanceReportCreated:
		data, err = notificationService.enrichMaintenanceReportData(notificationData)
		if err != nil {
			return nil, err
		}
	case enum.PanelReportCreated:
		data, err = notificationService.enrichPanelReportData(notificationData)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (notificationService *NotificationService) SendNotification(notification *entity.Notification, notificationType *entity.NotificationType) error {
	settings, _ := notificationService.notificationRepository.GetNotificationSettingByUserAndType(notificationService.db, notification.RecipientID, notification.TypeID)

	data, err := notificationService.dataCatcher(notificationType.Name, notification.Data)
	if err != nil {
		return err
	}

	if settings == nil {
		err = notificationService.CreateNotificationSettings(notification.RecipientID)
		if err != nil {
			return err
		}
	}

	if settings.IsPushEnabled {
		msg := notificationdto.PushNotificationResponse{
			ID:          notification.ID,
			Timestamp:   notification.CreatedAt,
			Type:        notificationType.Name.String(),
			Description: notificationType.Description,
			Data:        data,
			IsRead:      notification.IsRead,
			RecipientID: notification.RecipientID,
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		notificationService.wsHub.SendToUser(msg.RecipientID, websocket.MessageTypeNotification, msgBytes)
	}

	if settings.IsEmailEnabled {
		user, err := notificationService.userService.GetUserByID(notification.RecipientID)
		if err != nil {
			return err
		}
		if !user.EmailVerified {
			return nil
		}
		msg := struct {
			ToEmail      string      `json:"toEmail"`
			Subject      string      `json:"subject"`
			TemplateFile string      `json:"templateFile"`
			Data         interface{} `json:"data"`
		}{
			ToEmail:      user.Email,
			Subject:      notificationType.Name.String(),
			TemplateFile: notificationType.Name.EmailTemplatePath(),
			Data:         data,
		}

		if err := notificationService.emailService.SendEmail(msg.ToEmail, msg.Subject, msg.TemplateFile, msg.Data); err != nil {
			return err
		}
	}

	return nil
}

func (notificationService *NotificationService) MarkAsRead(notificationInfo notificationdto.NotificationInfoRequest) error {
	notification, err := notificationService.notificationRepository.GetNotificationByID(notificationService.db, notificationInfo.NotificationID)
	if err != nil {
		return err
	}
	if notification == nil {
		notFoundError := exception.NotFoundError{Item: notificationService.constants.Field.Notification}
		return notFoundError
	}
	if notification.RecipientID != notificationInfo.UserID {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: notificationService.constants.Field.Notification,
		}
		return forbiddenError
	}
	notification.IsRead = true

	err = notificationService.notificationRepository.UpdateNotification(notificationService.db, notification)
	if err != nil {
		return err
	}
	return nil
}

func (notificationService *NotificationService) GetNotificationsType() ([]notificationdto.NotificationTypeResponse, error) {
	notificationTypes, err := notificationService.notificationRepository.GetNotificationTypes(notificationService.db)
	if err != nil {
		return nil, err
	}
	notificationTypesResponse := make([]notificationdto.NotificationTypeResponse, len(notificationTypes))

	for i, notificationType := range notificationTypes {
		notificationTypesResponse[i] = notificationdto.NotificationTypeResponse{
			ID:            notificationType.ID,
			Name:          notificationType.Name.String(),
			Description:   notificationType.Description,
			SupportsEmail: notificationType.SupportsEmail,
			SupportsPush:  notificationType.SupportsPush,
		}
	}
	return notificationTypesResponse, nil
}

func (notificationService *NotificationService) GetUserNotifications(notificationsRequest notificationdto.NotificationListRequest) ([]notificationdto.NotificationListResponse, int64, error) {
	options := postgres.NewQueryOptions().
		WithPagination(notificationsRequest.Limit, notificationsRequest.Offset).
		WithSorting(notificationService.getSortByColumn(notificationsRequest.SortBy), notificationsRequest.Asc)

	notifications, err := notificationService.notificationRepository.GetNotificationsByTypesAndUserID(notificationService.db, notificationsRequest.UserID, notificationsRequest.Types, options)
	if err != nil {
		return nil, 0, err
	}
	notificationsResponse := make([]notificationdto.NotificationListResponse, len(notifications))

	for i, notification := range notifications {
		notificationType, err := notificationService.notificationRepository.GetNotificationTypeByID(notificationService.db, notification.TypeID)
		if err != nil {
			return nil, 0, err
		}
		if notificationType == nil {
			continue
		}

		data, _ := notificationService.dataCatcher(notificationType.Name, notification.Data)

		notificationTypeResponse := notificationdto.NotificationTypeResponse{
			ID:            notificationType.ID,
			Name:          notificationType.Name.String(),
			Description:   notificationType.Description,
			SupportsEmail: notificationType.SupportsEmail,
			SupportsPush:  notificationType.SupportsPush,
		}

		notificationsResponse[i] = notificationdto.NotificationListResponse{
			ID:     notification.ID,
			Type:   notificationTypeResponse,
			Data:   data,
			IsRead: notification.IsRead,
		}
	}

	count, err := notificationService.notificationRepository.CountNotificationsByTypesAndUserID(notificationService.db, notificationsRequest.UserID, notificationsRequest.Types)
	if err != nil {
		return nil, 0, err
	}

	return notificationsResponse, count, nil
}

func (notificationService *NotificationService) CreateNotificationSettings(userID uint) error {
	notificationTypes, err := notificationService.notificationRepository.GetNotificationTypes(notificationService.db)
	if err != nil {
		return err
	}
	for _, notificationType := range notificationTypes {
		setting, err := notificationService.notificationRepository.GetNotificationSettingByUserAndType(notificationService.db, userID, notificationType.ID)
		if err != nil {
			return err
		}
		if setting != nil {
			continue
		}

		setting = &entity.NotificationSetting{
			UserID:         userID,
			TypeID:         notificationType.ID,
			IsEmailEnabled: notificationType.SupportsEmail,
			IsPushEnabled:  notificationType.SupportsPush,
		}
		err = notificationService.notificationRepository.CreateNotificationSetting(notificationService.db, setting)
		if err != nil {
			return err
		}
	}
	return nil
}

func (notificationService *NotificationService) GetUserNotificationSettings(userID uint) ([]notificationdto.NotificationSettingResponse, error) {
	settings, err := notificationService.notificationRepository.GetNotificationSettingByUserID(notificationService.db, userID)
	if err != nil {
		return nil, err
	}
	settingsResponse := make([]notificationdto.NotificationSettingResponse, len(settings))

	for i, setting := range settings {
		notificationType, err := notificationService.notificationRepository.GetNotificationTypeByID(notificationService.db, setting.TypeID)
		if err != nil {
			return nil, err
		}
		if notificationType == nil {
			continue
		}
		notificationTypeResponse := notificationdto.NotificationTypeResponse{
			ID:            notificationType.ID,
			Name:          notificationType.Name.String(),
			Description:   notificationType.Description,
			SupportsEmail: notificationType.SupportsEmail,
			SupportsPush:  notificationType.SupportsPush,
		}
		settingsResponse[i] = notificationdto.NotificationSettingResponse{
			ID:               setting.ID,
			NotificationType: notificationTypeResponse,
			IsEmailEnabled:   setting.IsEmailEnabled,
			IsPushEnabled:    setting.IsPushEnabled,
		}
	}
	return settingsResponse, nil
}

func (notificationService *NotificationService) UpdateNotificationSettings(newSettingInfo notificationdto.UpdateSettingsRequest) error {
	setting, err := notificationService.notificationRepository.GetNotificationSettingByID(notificationService.db, newSettingInfo.SettingID)
	if err != nil {
		return err
	}
	if setting == nil {
		notFoundError := exception.NotFoundError{Item: notificationService.constants.Field.NotificationSetting}
		return notFoundError
	}
	notificationType, err := notificationService.notificationRepository.GetNotificationTypeByID(notificationService.db, setting.TypeID)
	if err != nil {
		return err
	}
	if notificationType == nil {
		notFoundError := exception.NotFoundError{Item: notificationService.constants.Field.NotificationType}
		return notFoundError
	}
	setting.IsEmailEnabled = newSettingInfo.IsEmailEnabled && notificationType.SupportsEmail
	setting.IsPushEnabled = newSettingInfo.IsPushEnabled && notificationType.SupportsPush
	err = notificationService.notificationRepository.UpdateNotificationSetting(notificationService.db, setting)
	if err != nil {
		return err
	}
	return nil
}
