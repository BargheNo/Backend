package enum

type NotificationType uint

const (
	ChatNotificationType NotificationType = iota + 1
	CorpSendBidNotificationType
	ReportCreated
)

var notificationDescriptions = map[NotificationType]string{
	ChatNotificationType:        "پیام جدید",
	CorpSendBidNotificationType: "پیشنهاد جدید از سوی شرکت",
	ReportCreated:               "گزارش جدید",
}

var supportsEmail = map[NotificationType]bool{
	ChatNotificationType:        false,
	CorpSendBidNotificationType: true,
	ReportCreated:               true,
}

var supportsPush = map[NotificationType]bool{
	ChatNotificationType:        true,
	CorpSendBidNotificationType: true,
	ReportCreated:               true,
}

var notificationEmailTemplate = map[NotificationType]string{
	ChatNotificationType:        "",
	CorpSendBidNotificationType: "/sample/sample.html",
	ReportCreated:               "/sample/sample.html",
}

func (notificationType NotificationType) String() string {
	switch notificationType {
	case ChatNotificationType:
		return "پیام جدید"
	case CorpSendBidNotificationType:
		return "پیشنهادات قیمت جدید"
	case ReportCreated:
		return "گزارشات جدید"
	}
	return ""
}

func (notificationType NotificationType) Description() string {
	if description, ok := notificationDescriptions[notificationType]; ok {
		return description
	}
	return ""
}

func (notificationType NotificationType) EmailTemplatePath() string {
	if template, ok := notificationEmailTemplate[notificationType]; ok {
		return template
	}
	return ""
}

func (notificationType NotificationType) SupportsEmail() bool {
	if support, ok := supportsEmail[notificationType]; ok {
		return support
	}
	return false
}

func (notificationType NotificationType) SupportsPush() bool {
	if support, ok := supportsPush[notificationType]; ok {
		return support
	}
	return true
}

func GetAllNotificationTypes() []NotificationType {
	return []NotificationType{
		ChatNotificationType,
		CorpSendBidNotificationType,
		ReportCreated,
	}
}
