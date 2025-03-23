package enum

type BidStatus uint

const (
	Pending BidStatus = iota + 1
	Accepted
	// Expired
)

func (s BidStatus) String() string {
	switch s {
	case Pending:
		return "pending"
	case Accepted:
		return "accepted"
		// case Expired:
		// return "expired"
	}
	return "unknown"
}

func GetAllBidStatuses() []BidStatus {
	return []BidStatus{
		Pending,
		Accepted,
		// Expired,
	}
}
