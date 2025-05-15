package enum

type NotificationType uint

const (
	ChatNotificationType NotificationType = iota + 1
	CorpSendBidNotificationType
	PanelReportCreated
	MaintenanceReportCreated
)

var notificationDescriptions = map[NotificationType]string{
	ChatNotificationType:        "پیام جدید",
	CorpSendBidNotificationType: "پیشنهاد جدید از سوی شرکت",
	PanelReportCreated:          "گزارش جدید از پنل",
	MaintenanceReportCreated:    "گزارش جدید از تعمیرات",
}

var supportsEmail = map[NotificationType]bool{
	ChatNotificationType:        false,
	CorpSendBidNotificationType: true,
	PanelReportCreated:          true,
	MaintenanceReportCreated:    true,
}

var supportsPush = map[NotificationType]bool{
	ChatNotificationType:        true,
	CorpSendBidNotificationType: true,
	PanelReportCreated:          true,
	MaintenanceReportCreated:    true,
}

var notificationEmailTemplate = map[NotificationType]string{
	ChatNotificationType:        "",
	CorpSendBidNotificationType: "/get_bid/fa.html",
	PanelReportCreated:          "/panel_report/fa.html",
	MaintenanceReportCreated:    "/maintenance_report/fa.html",
}

func (notificationType NotificationType) String() string {
	switch notificationType {
	case ChatNotificationType:
		return "پیام جدید"
	case CorpSendBidNotificationType:
		return "درخواست نصب پنل"
	case PanelReportCreated:
		return "گزارشات جدید از پنل"
	case MaintenanceReportCreated:
		return "گزارشات جدید از تعمیرات"
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
		PanelReportCreated,
		MaintenanceReportCreated,
	}
}
