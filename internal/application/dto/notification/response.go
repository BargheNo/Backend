package notificationdto

type NotificationListResponse struct {
	ID             uint                     `json:"id"`
	Type           NotificationTypeResponse `json:"type"`
	Message        string                   `json:"message"`
	AdditionalData string                   `json:"additionalData"`
	IsRead         bool                     `json:"isRead"`
}

type NotificationTypeResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type NotificationSettingResponse struct {
	UserID           uint                     `gorm:"not null;index"`
	TypeID           uint                     `gorm:"not null;index"`
	NotificationType NotificationTypeResponse `gorm:"foreignKey:TypeID"`
	IsEmailEnabled   bool                     `gorm:"default:true"`
	IsPushEnabled    bool                     `gorm:"default:true"`
}
