package sortby

type MaintenanceSortBy uint

const (
	MaintenanceSortByCreatedAt MaintenanceSortBy = iota + 1
)

func (s MaintenanceSortBy) Name() string {
	switch s {
	case MaintenanceSortByCreatedAt:
		return "تاریخ ایجاد"
	}
	return "ناشناس"
}

func (s MaintenanceSortBy) DBColumn() string {
	switch s {
	case MaintenanceSortByCreatedAt:
		return "created_at"
	}
	return ""
}

func GetMaintenanceSortableColumns() map[MaintenanceSortBy]bool {
	return map[MaintenanceSortBy]bool{
		MaintenanceSortByCreatedAt: true,
	}
}
