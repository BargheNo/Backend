package sortby

type UserSortBy uint

const (
	UserSortByCreatedAt UserSortBy = iota + 1
	UserSortByName
)

func (s UserSortBy) Name() string {
	switch s {
	case UserSortByCreatedAt:
		return "تاریخ ایجاد"
	case UserSortByName:
		return "نام"
	}
	return "ناشناس"
}

func (s UserSortBy) DBColumn() string {
	switch s {
	case UserSortByCreatedAt:
		return "created_at"
	case UserSortByName:
		return "last_name"
	}
	return ""
}

func GetUserSortableColumns() map[UserSortBy]bool {
	return map[UserSortBy]bool{
		UserSortByCreatedAt: true,
		UserSortByName:      true,
	}
}
