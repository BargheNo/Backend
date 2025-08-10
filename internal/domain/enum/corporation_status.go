package enum

type CorporationStatus uint

const (
	CorpStatusApproved CorporationStatus = iota + 1
	CorpStatusAwaitingApproval
	CorpStatusSuspend
	CorpStatusRejected
	CorpStatusAll
)

func (s CorporationStatus) String() string {
	switch s {
	case CorpStatusApproved:
		return "تایید شده"
	case CorpStatusAwaitingApproval:
		return "در انتظار تایید"
	case CorpStatusSuspend:
		return "معلق"
	case CorpStatusRejected:
		return "رد شده"
	case CorpStatusAll:
		return "همه"
	}
	return "unknown"
}

func GetAllCorporationStatuses() []CorporationStatus {
	return []CorporationStatus{
		CorpStatusApproved,
		CorpStatusAwaitingApproval,
		CorpStatusSuspend,
		CorpStatusRejected,
		CorpStatusAll,
	}
}
