package sortby

type ChatSortBy uint

const (
	ChatSortByCreatedAt ChatSortBy = iota + 1
)

func (s ChatSortBy) Name() string {
	switch s {
	case ChatSortByCreatedAt:
		return "تاریخ ایجاد"
	}
	return "ناشناس"
}

func (s ChatSortBy) DBColumn() string {
	switch s {
	case ChatSortByCreatedAt:
		return "created_at"
	}
	return ""
}

func GetChatSortableColumns() map[ChatSortBy]bool {
	return map[ChatSortBy]bool{
		ChatSortByCreatedAt: true,
	}
}
