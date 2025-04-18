package enum

type ReportStatus uint

const (
	ReportStatusPending ReportStatus = iota + 1
	ReportStatusResolved
)

func (s ReportStatus) String() string {
	switch s {
	case ReportStatusPending:
		return "pending"
	case ReportStatusResolved:
		return "resolved"
	default:
		return "unknown"
	}
}

func GetAllReportStatuses() []ReportStatus {
	return []ReportStatus{
		ReportStatusPending,
		ReportStatusResolved,
	}
}
