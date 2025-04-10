package serviceimpl

import (
	"encoding/json"
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	notificationdto "github.com/BargheNo/Backend/internal/application/dto/notification"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/BargheNo/Backend/internal/infrastructure/websocket"
)

type NotificationService struct {
	constants              *bootstrap.Constants
	userService            service.UserService
	notificationRepository repository.NotificationRepository
	wsHub                  *websocket.Hub
	db                     database.Database
}

func NewNotificationService(
	constants *bootstrap.Constants,
	userService service.UserService,
	notificationRepository repository.NotificationRepository,
	wsHub *websocket.Hub,
	db database.Database,
) *NotificationService {
	return &NotificationService{
		constants:              constants,
		userService:            userService,
		notificationRepository: notificationRepository,
		wsHub:                  wsHub,
		db:                     db,
	}
}

func (notificationService *NotificationService) CreateNotification(
	typeID, recipientID uint,
	additionalData map[string]string,
) {
	notificationService.userService.GetUserCredential(recipientID)
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
		panic(err)
	}

	settings, exist := notificationService.notificationRepository.GetNotificationSettingByUserAndType(notificationService.db, recipientID, typeID)
	if !exist {
		return
	}

	if settings.IsEmailEnabled {
		if err := notificationService.SendNotification(notification); err != nil {
			panic(err)
		}
	}
}

func (notificationService *NotificationService) SendNotification(notification *entity.Notification) error {
	notificationType, exist := notificationService.notificationRepository.GetNotificationTypeByID(notificationService.db, notification.TypeID)
	if !exist {
		return exception.NotFoundError{Item: notificationService.constants.Field.NotificationType}
	}

	var additionalData map[string]string
	if err := json.Unmarshal([]byte(notification.AdditionalData), &additionalData); err != nil {
		additionalData = nil
	}

	payload := websocket.NotificationPayload{
		ID:             notification.ID,
		Title:          notificationType.Name,
		Description:    notificationType.Description,
		AdditionalData: additionalData,
		Type:           notificationType.Name,
		IsRead:         notification.IsRead,
		CreatedAt:      notification.CreatedAt.Format(time.RFC3339),
	}

	content, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	notificationService.wsHub.SendToUser(notification.RecipientID, websocket.MessageTypeNotification, content)

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
	notificationService.userService.GetUserCredential(userID)

	notifications := notificationService.notificationRepository.GetNotificationsByUserID(notificationService.db, userID)
	notificationsResponse := make([]notificationdto.NotificationListResponse, len(notifications))

	for i, notification := range notifications {
		notificationType, exist := notificationService.notificationRepository.GetNotificationTypeByID(notificationService.db, notification.TypeID)
		if !exist {
			continue
		}
		notificationsResponse[i] = notificationdto.NotificationListResponse{
			ID:             notification.ID,
			Type:           notificationdto.NotificationTypeResponse{Name: notificationType.Name, Description: notificationType.Description},
			AdditionalData: notification.AdditionalData,
			IsRead:         notification.IsRead,
		}
	}
	return notificationsResponse
}

func (notificationService *NotificationService) CreateNotificationSettings(userID uint) {
	notificationService.userService.GetUserCredential(userID)
	notificationTypes := notificationService.notificationRepository.GetNotificationTypes(notificationService.db)
	for _, notificationType := range notificationTypes {
		setting := &entity.NotificationSetting{
			UserID:         userID,
			TypeID:         notificationType.ID,
			IsEmailEnabled: true,
			IsPushEnabled:  true,
		}
		err := notificationService.notificationRepository.CreateNotificationSetting(notificationService.db, setting)
		if err != nil {
			panic(err)
		}
	}
}

func (notificationService *NotificationService) GetUserNotificationSettings(userID uint) []notificationdto.NotificationSettingResponse {
	notificationService.userService.GetUserCredential(userID)

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
			NotificationType: notificationdto.NotificationTypeResponse{Name: notificationType.Name, Description: notificationType.Description},
			IsEmailEnabled:   setting.IsEmailEnabled,
			IsPushEnabled:    setting.IsPushEnabled,
		}
	}
	return settingsResponse
}

func (notificationService *NotificationService) UpdateNotificationSettings(newSettingInfo notificationdto.UpdateSettingsRequest) {
	notificationService.userService.GetUserCredential(newSettingInfo.UserID)
	setting, exist := notificationService.notificationRepository.GetNotificationSettingByID(notificationService.db, newSettingInfo.SettingID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: notificationService.constants.Field.NotificationType}
		panic(notFoundError)
	}
	setting.IsEmailEnabled = newSettingInfo.IsEmailEnabled
	setting.IsPushEnabled = newSettingInfo.IsPushEnabled
	notificationService.notificationRepository.UpdateNotificationSetting(notificationService.db, setting)
}
