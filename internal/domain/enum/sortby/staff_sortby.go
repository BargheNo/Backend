package sortby

type StaffSortBy uint

const (
	StaffSortByCreatedAt StaffSortBy = iota + 1
	StaffSortByName
)

func (s StaffSortBy) Name() string {
	switch s {
	case StaffSortByCreatedAt:
		return "تاریخ ایجاد"
	case StaffSortByName:
		return "نام"
	}
	return "ناشناس"
}

func (s StaffSortBy) DBColumn() string {
	switch s {
	case StaffSortByCreatedAt:
		return "corporation_staffs.created_at"
	case StaffSortByName:
		return "users.last_name"
	}
	return ""
}

func GetStaffSortableColumns() map[StaffSortBy]bool {
	return map[StaffSortBy]bool{
		StaffSortByCreatedAt: true,
		StaffSortByName:      true,
	}
}
