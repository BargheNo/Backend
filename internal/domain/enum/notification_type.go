package enum

type NotificationType uint

const (
	ChatNotificationType NotificationType = iota + 1
	CorpSendBidNotificationType
)

var notificationDescriptions = map[NotificationType]string{
	ChatNotificationType:        "پیام جدید",
	CorpSendBidNotificationType: "پیشنهاد جدید از سوی شرکت",
}

var notificationEmailTemplate = map[NotificationType]string{
	ChatNotificationType:        "/sample/sample.html",
	CorpSendBidNotificationType: "/sample/sample.html",
}

func (NotificationType NotificationType) String() string {
	switch NotificationType {
	case ChatNotificationType:
		return "chat"
	case CorpSendBidNotificationType:
		return "bid"
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
		CorpSendBidNotificationType,
	}
}
