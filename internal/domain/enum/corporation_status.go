package enum

type CorporationStatus uint

const (
	CorpStatusApproved CorporationStatus = iota + 1
	CorpStatusAwaitingApproval
	CorpStatusRejected
)

func (s CorporationStatus) String() string {
	switch s {
	case CorpStatusApproved:
		return "approved"
	case CorpStatusAwaitingApproval:
		return "awaiting_approval"
	case CorpStatusRejected:
		return "rejected"
	}
	return "unknown"
}

func GetAllCorporationStatuses() []CorporationStatus {
	return []CorporationStatus{
		CorpStatusApproved,
		CorpStatusAwaitingApproval,
		CorpStatusRejected,
	}
}
