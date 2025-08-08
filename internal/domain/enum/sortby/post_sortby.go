package sortby

type PostSortBy uint

const (
	PostSortByCreatedAt PostSortBy = iota + 1
	PostSortByLikeCount
)

func (s PostSortBy) Name() string {
	switch s {
	case PostSortByCreatedAt:
		return "تاریخ ایجاد"
	case PostSortByLikeCount:
		return "تعداد لایک"
	}
	return "ناشناس"
}

func (s PostSortBy) DBColumn() string {
	switch s {
	case PostSortByCreatedAt:
		return "created_at"
	case PostSortByLikeCount:
		return "like_count"
	}
	return ""
}

func GetPostSortableColumns() map[PostSortBy]bool {
	return map[PostSortBy]bool{
		PostSortByCreatedAt: true,
		PostSortByLikeCount: true,
	}
}
