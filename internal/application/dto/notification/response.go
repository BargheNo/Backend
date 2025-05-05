package notificationdto

type NotificationListResponse struct {
	ID             uint                     `json:"id"`
	Type           NotificationTypeResponse `json:"type"`
	AdditionalData string                   `json:"additionalData"`
	IsRead         bool                     `json:"isRead"`
}

type NotificationTypeResponse struct {
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
