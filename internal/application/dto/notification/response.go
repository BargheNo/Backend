package notificationdto

import (
	"time"
)

type NotificationListResponse struct {
	ID     uint                     `json:"id"`
	Type   NotificationTypeResponse `json:"type"`
	Data   map[string]interface{}   `json:"data"`
	IsRead bool                     `json:"isRead"`
}

type NotificationTypeResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	SupportsEmail bool   `json:"supportsEmail"`
	SupportsPush  bool   `json:"supportsPush"`
}

type NotificationSettingResponse struct {
	ID               uint                     `json:"id"`
	NotificationType NotificationTypeResponse `json:"notificationType"`
	IsEmailEnabled   bool                     `json:"isEmailEnabled"`
	IsPushEnabled    bool                     `json:"isPushEnabled"`
}

type PushNotificationResponse struct {
	ID          uint                   `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data"`
	IsRead      bool                   `json:"isRead"`
	RecipientID uint                   `json:"recipientID"`
}
