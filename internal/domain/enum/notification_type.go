package enum

type NotificationType uint

const (
	ChatNotificationType NotificationType = iota + 1
)

func (NotificationType NotificationType) String() string {
	switch NotificationType {
	case ChatNotificationType:
		return "chat"
	}
	return ""
}

func GetAllNotificationTypes() []NotificationType {
	return []NotificationType{
		ChatNotificationType,
	}
}
