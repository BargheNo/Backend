package serviceimpl

import (
	"encoding/json"

	"github.com/BargheNo/Backend/bootstrap"
	biddto "github.com/BargheNo/Backend/internal/application/dto/bid"
	notificationdto "github.com/BargheNo/Backend/internal/application/dto/notification"
	reportdto "github.com/BargheNo/Backend/internal/application/dto/report"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/message"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/websocket"
)

type NotificationService struct {
	constants              *bootstrap.Constants
	userService            service.UserService
	emailService           service.EmailService
	bidService             service.BidService
	reportService          service.ReportService
	notificationRepository repository.NotificationRepository
	wsHub                  *websocket.Hub
	rabbitMQ               message.Broker
	db                     database.Database
}

type NotificationServiceDeps struct {
	Constants              *bootstrap.Constants
	UserService            service.UserService
	EmailService           service.EmailService
	BidService             service.BidService
	ReportService          service.ReportService
	NotificationRepository repository.NotificationRepository
	WSHub                  *websocket.Hub
	RabbitMQ               message.Broker
	DB                     database.Database
}

func NewNotificationService(deps NotificationServiceDeps) *NotificationService {
	return &NotificationService{
		constants:              deps.Constants,
		userService:            deps.UserService,
		emailService:           deps.EmailService,
		bidService:             deps.BidService,
		reportService:          deps.ReportService,
		notificationRepository: deps.NotificationRepository,
		wsHub:                  deps.WSHub,
		rabbitMQ:               deps.RabbitMQ,
		db:                     deps.DB,
	}
}

func (notificationService *NotificationService) CreateAndSendNotification(typeName enum.NotificationType, recipientID uint, data []byte) error {
	notificationService.userService.DoesUserExist(recipientID)
	notificationType, exist := notificationService.notificationRepository.GetNotificationTypeByName(notificationService.db, typeName)
	if !exist {
		notFoundError := exception.NotFoundError{Item: notificationService.constants.Field.NotificationType}
		panic(notFoundError)
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
	bid := notificationService.bidService.GetRequestAnonymousBid(requestInfo)

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
	report := notificationService.reportService.GetMaintenanceReport(reportData.ReportID)

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
	report := notificationService.reportService.GetPanelReport(reportData.ReportID)

	reportBytes, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(reportBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (notificationService *NotificationService) SendNotification(notification *entity.Notification, notificationType *entity.NotificationType) error {
	settings, _ := notificationService.notificationRepository.GetNotificationSettingByUserAndType(notificationService.db, notification.RecipientID, notification.TypeID)

	var data map[string]interface{}
	var err error

	switch notificationType.Name {
	case enum.CorpSendBidNotificationType:
		data, err = notificationService.enrichBidData(notification.Data)
		if err != nil {
			return err
		}
	case enum.MaintenanceReportCreated:
		data, err = notificationService.enrichMaintenanceReportData(notification.Data)
		if err != nil {
			return err
		}
	case enum.PanelReportCreated:
		data, err = notificationService.enrichPanelReportData(notification.Data)
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
		user := notificationService.userService.GetUserByID(notification.RecipientID)
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

func (notificationService *NotificationService) MarkAsRead(notificationInfo notificationdto.NotificationInfoRequest) {
	notification, exist := notificationService.notificationRepository.GetNotificationByID(notificationService.db, notificationInfo.NotificationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: notificationService.constants.Field.NotificationType}
		panic(notFoundError)
	}
	if notification.RecipientID != notificationInfo.UserID {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: notificationService.constants.Field.Notification,
		}
		panic(forbiddenError)
	}
	notification.IsRead = true

	err := notificationService.notificationRepository.UpdateNotification(notificationService.db, notification)
	if err != nil {
		panic(err)
	}
}

func (notificationService *NotificationService) GetNotificationsType() []notificationdto.NotificationTypeResponse {
	notificationTypes := notificationService.notificationRepository.GetNotificationTypes(notificationService.db)
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
	return notificationTypesResponse
}

func (notificationService *NotificationService) GetUserNotifications(notificationsRequest notificationdto.NotificationListRequest) []notificationdto.NotificationListResponse {
	notificationService.userService.DoesUserExist(notificationsRequest.UserID)
	paginationModifier := repositoryimpl.NewPaginationModifier(notificationsRequest.Limit, notificationsRequest.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	notifications := notificationService.notificationRepository.GetNotificationsByTypesAndUserID(notificationService.db, notificationsRequest.UserID, notificationsRequest.Types, paginationModifier, sortingModifier)
	notificationsResponse := make([]notificationdto.NotificationListResponse, len(notifications))

	for i, notification := range notifications {
		notificationType, exist := notificationService.notificationRepository.GetNotificationTypeByID(notificationService.db, notification.TypeID)
		if !exist {
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal(notification.Data, &data); err != nil {
			continue
		}

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
	return notificationsResponse
}

func (notificationService *NotificationService) CreateNotificationSettings(userID uint) {
	notificationService.userService.DoesUserExist(userID)
	notificationTypes := notificationService.notificationRepository.GetNotificationTypes(notificationService.db)
	for _, notificationType := range notificationTypes {
		_, exist := notificationService.notificationRepository.GetNotificationSettingByUserAndType(notificationService.db, userID, notificationType.ID)
		if exist {
			continue
		}
		setting := &entity.NotificationSetting{
			UserID:         userID,
			TypeID:         notificationType.ID,
			IsEmailEnabled: notificationType.SupportsEmail,
			IsPushEnabled:  notificationType.SupportsPush,
		}
		err := notificationService.notificationRepository.CreateNotificationSetting(notificationService.db, setting)
		if err != nil {
			panic(err)
		}
	}
}

func (notificationService *NotificationService) GetUserNotificationSettings(userID uint) []notificationdto.NotificationSettingResponse {
	notificationService.userService.DoesUserExist(userID)

	settings := notificationService.notificationRepository.GetNotificationSettingByUserID(notificationService.db, userID)
	settingsResponse := make([]notificationdto.NotificationSettingResponse, len(settings))

	for i, setting := range settings {
		notificationType, exist := notificationService.notificationRepository.GetNotificationTypeByID(notificationService.db, setting.TypeID)
		if !exist {
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
	return settingsResponse
}

func (notificationService *NotificationService) UpdateNotificationSettings(newSettingInfo notificationdto.UpdateSettingsRequest) {
	notificationService.userService.DoesUserExist(newSettingInfo.UserID)
	setting, exist := notificationService.notificationRepository.GetNotificationSettingByID(notificationService.db, newSettingInfo.SettingID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: notificationService.constants.Field.NotificationSetting}
		panic(notFoundError)
	}
	notificationType, exist := notificationService.notificationRepository.GetNotificationTypeByID(notificationService.db, setting.TypeID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: notificationService.constants.Field.NotificationType}
		panic(notFoundError)
	}
	setting.IsEmailEnabled = newSettingInfo.IsEmailEnabled && notificationType.SupportsEmail
	setting.IsPushEnabled = newSettingInfo.IsPushEnabled && notificationType.SupportsPush
	notificationService.notificationRepository.UpdateNotificationSetting(notificationService.db, setting)
}
