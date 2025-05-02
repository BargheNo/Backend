package serviceimpl

import (
	"encoding/json"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	notificationdto "github.com/BargheNo/Backend/internal/application/dto/notification"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/message"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/BargheNo/Backend/internal/infrastructure/websocket"
)

type NotificationService struct {
	constants              *bootstrap.Constants
	userService            service.UserService
	notificationRepository repository.NotificationRepository
	wsHub                  *websocket.Hub
	rabbitMQ               message.Broker
	db                     database.Database
}

func NewNotificationService(
	constants *bootstrap.Constants,
	userService service.UserService,
	notificationRepository repository.NotificationRepository,
	wsHub *websocket.Hub,
	rabbitMQ message.Broker,
	db database.Database,
) *NotificationService {
	return &NotificationService{
		constants:              constants,
		userService:            userService,
		notificationRepository: notificationRepository,
		wsHub:                  wsHub,
		rabbitMQ:               rabbitMQ,
		db:                     db,
	}
}

func (notificationService *NotificationService) CreateAndSendNotification(
	typeID, recipientID uint, additionalData interface{}) error {
	notificationService.userService.DoesUserExist(recipientID)
	additionalDataJSON, err := json.Marshal(additionalData)
	if err != nil {
		panic(err)
	}

	notification := &entity.Notification{
		TypeID:         typeID,
		RecipientID:    recipientID,
		AdditionalData: string(additionalDataJSON),
		IsRead:         false,
	}

	err = notificationService.notificationRepository.CreateNotification(notificationService.db, notification)
	if err != nil {
		return err
	}

	settings, _ := notificationService.notificationRepository.GetNotificationSettingByUserAndType(notificationService.db, recipientID, typeID)

	if settings.IsEmailEnabled {
		if err := notificationService.SendNotification(notification); err != nil {
			return err
		}
	}

	return nil
}

func (notificationService *NotificationService) SendNotification(notification *entity.Notification) error {
	settings, _ := notificationService.notificationRepository.GetNotificationSettingByUserAndType(notificationService.db, notification.RecipientID, notification.TypeID)
	notificationType, _ := notificationService.notificationRepository.GetNotificationTypeByID(notificationService.db, settings.TypeID)
	if settings.IsPushEnabled {
		msg := struct {
			ID             uint      `json:"id"`
			Timestamp      time.Time `json:"timestamp"`
			Type           string    `json:"type"`
			Description    string    `json:"description"`
			AdditionalData string    `json:"additionalData"`
			IsRead         bool      `json:"isRead"`
			RecipientID    uint      `json:"recipientID"`
		}{
			ID:             notification.ID,
			Timestamp:      notification.CreatedAt,
			Type:           notificationType.Name.String(),
			Description:    notificationType.Description,
			AdditionalData: notification.AdditionalData,
			IsRead:         notification.IsRead,
			RecipientID:    notification.RecipientID,
		}
		if err := notificationService.rabbitMQ.PublishMessage(notificationService.constants.RabbitMQ.Events.NotificationsPush, msg); err != nil {
			return err
		}
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
			Subject:      "hello",
			TemplateFile: "/sample/sample.html",
			Data:         nil,
		}
		if err := notificationService.rabbitMQ.PublishMessage(notificationService.constants.RabbitMQ.Events.NotificationsEmail, msg); err != nil {
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

func (notificationService *NotificationService) GetUserNotifications(userID uint) []notificationdto.NotificationListResponse {
	notificationService.userService.DoesUserExist(userID)

	notifications := notificationService.notificationRepository.GetNotificationsByUserID(notificationService.db, userID)
	notificationsResponse := make([]notificationdto.NotificationListResponse, len(notifications))

	for i, notification := range notifications {
		notificationType, exist := notificationService.notificationRepository.GetNotificationTypeByID(notificationService.db, notification.TypeID)
		if !exist {
			continue
		}
		notificationsResponse[i] = notificationdto.NotificationListResponse{
			ID:             notification.ID,
			Type:           notificationdto.NotificationTypeResponse{Name: notificationType.Name.String(), Description: notificationType.Description},
			AdditionalData: notification.AdditionalData,
			IsRead:         notification.IsRead,
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
			IsEmailEnabled: true,
			IsPushEnabled:  false,
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
		settingsResponse[i] = notificationdto.NotificationSettingResponse{
			UserID:           userID,
			TypeID:           setting.TypeID,
			NotificationType: notificationdto.NotificationTypeResponse{Name: notificationType.Name.String(), Description: notificationType.Description},
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
		notFoundError := exception.NotFoundError{Item: notificationService.constants.Field.NotificationType}
		panic(notFoundError)
	}
	setting.IsEmailEnabled = newSettingInfo.IsEmailEnabled
	setting.IsPushEnabled = newSettingInfo.IsPushEnabled
	notificationService.notificationRepository.UpdateNotificationSetting(notificationService.db, setting)
}
