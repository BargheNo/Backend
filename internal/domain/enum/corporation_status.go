package enum

type CorporationStatus uint

const (
	Approved CorporationStatus = iota + 1
	AwaitingApproval
	Rejected
)

func (s CorporationStatus) String() string {
	switch s {
	case Approved:
		return "approved"
	case AwaitingApproval:
		return "awaiting_approval"
	case Rejected:
		return "rejected"
	}
	return "unknown"
}

func GetAllCorporationStatuses() []CorporationStatus {
	return []CorporationStatus{
		Approved,
		AwaitingApproval,
		Rejected,
	}
}
