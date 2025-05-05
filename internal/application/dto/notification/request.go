package notificationdto

type NotificationInfoRequest struct {
	NotificationID uint
	UserID         uint
}

type NotificationListRequest struct {
	Types  []uint
	UserID uint
}

type UpdateSettingsRequest struct {
	SettingID      uint
	UserID         uint
	IsEmailEnabled bool
	IsPushEnabled  bool
}
