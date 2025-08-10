package sortby

type ReportSortBy uint

const (
	ReportSortByCreatedAt ReportSortBy = iota + 1
)

func (s ReportSortBy) Name() string {
	switch s {
	case ReportSortByCreatedAt:
		return "تاریخ ایجاد"
	}
	return "ناشناس"
}

func (s ReportSortBy) DBColumn() string {
	switch s {
	case ReportSortByCreatedAt:
		return "created_at"
	}
	return ""
}

func GetReportSortableColumns() map[ReportSortBy]bool {
	return map[ReportSortBy]bool{
		ReportSortByCreatedAt: true,
	}
}
