package enum

type NotificationType uint

const (
	ChatNotificationType NotificationType = iota + 1
)

var notificationDescriptions = map[NotificationType]string{
	ChatNotificationType: "پیام جدید",
}

var notificationEmailTemplate = map[NotificationType]string{
	ChatNotificationType: "/sample/sample.html",
}

func (NotificationType NotificationType) String() string {
	switch NotificationType {
	case ChatNotificationType:
		return "chat"
	}
	return ""
}

func (NotificationType NotificationType) Description() string {
	if description, ok := notificationDescriptions[NotificationType]; ok {
		return description
	}
	return ""
}

func (NotificationType NotificationType) EmailTemplatePath() string {
	if description, ok := notificationEmailTemplate[NotificationType]; ok {
		return description
	}
	return ""
}

func GetAllNotificationTypes() []NotificationType {
	return []NotificationType{
		ChatNotificationType,
	}
}
