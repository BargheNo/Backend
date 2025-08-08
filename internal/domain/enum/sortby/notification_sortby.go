package sortby

type NotificationSortBy uint

const (
	NotificationSortByCreatedAt NotificationSortBy = iota + 1
)

func (s NotificationSortBy) Name() string {
	switch s {
	case NotificationSortByCreatedAt:
		return "تاریخ ایجاد"
	}
	return "ناشناس"
}

func (s NotificationSortBy) DBColumn() string {
	switch s {
	case NotificationSortByCreatedAt:
		return "created_at"
	}
	return ""
}

func GetNotificationSortableColumns() map[NotificationSortBy]bool {
	return map[NotificationSortBy]bool{
		NotificationSortByCreatedAt: true,
	}
}
