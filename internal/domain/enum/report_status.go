package enum

type ReportStatus uint

const (
	ReportStatusPending ReportStatus = iota + 1
	ReportStatusResolved
	ReportStatusAll
)

func (s ReportStatus) String() string {
	switch s {
	case ReportStatusPending:
		return "درحال بررسی"
	case ReportStatusResolved:
		return "بررسی شده"
	case ReportStatusAll:
		return "همه"
	}
	return "unknown"
}

func GetAllReportStatuses() []ReportStatus {
	return []ReportStatus{
		ReportStatusPending,
		ReportStatusResolved,
		ReportStatusAll,
	}
}
