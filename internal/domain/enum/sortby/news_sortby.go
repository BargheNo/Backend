package sortby

type NewsSortBy uint

const (
	NewsSortByCreatedAt NewsSortBy = iota + 1
	NewsSortByLikeCount
)

func (s NewsSortBy) Name() string {
	switch s {
	case NewsSortByCreatedAt:
		return "تاریخ ایجاد"
	case NewsSortByLikeCount:
		return "تعداد لایک"
	}
	return "ناشناس"
}

func (s NewsSortBy) DBColumn() string {
	switch s {
	case NewsSortByCreatedAt:
		return "created_at"
	case NewsSortByLikeCount:
		return "like_count"
	}
	return ""
}

func GetNewsSortableColumns() map[NewsSortBy]bool {
	return map[NewsSortBy]bool{
		NewsSortByCreatedAt: true,
		NewsSortByLikeCount: true,
	}
}
