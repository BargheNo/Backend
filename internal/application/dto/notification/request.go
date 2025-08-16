package notificationdto

type NotificationInfoRequest struct {
	NotificationID uint
	UserID         uint
}

type NotificationListRequest struct {
	Types  []uint
	IsRead bool
	UserID uint
	Offset int
	Limit  int
	SortBy uint
	Asc    bool
}

type UpdateSettingsRequest struct {
	SettingID      uint
	UserID         uint
	IsEmailEnabled bool
	IsPushEnabled  bool
}
