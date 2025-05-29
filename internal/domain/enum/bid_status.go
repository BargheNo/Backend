package enum

type BidStatus uint

const (
	BidStatusPending BidStatus = iota + 1
	BidStatusAccepted
	BidStatusExpired
	BidStatusRejected
	BidStatusCanceled
	BidStatusAll
)

func (s BidStatus) String() string {
	switch s {
	case BidStatusPending:
		return "pending"
	case BidStatusAccepted:
		return "accepted"
	case BidStatusExpired:
		return "expired"
	case BidStatusRejected:
		return "rejected"
	case BidStatusCanceled:
		return "canceled"
	}
	return "unknown"
}

func GetAllBidStatuses() []BidStatus {
	return []BidStatus{
		BidStatusPending,
		BidStatusAccepted,
		BidStatusExpired,
		BidStatusRejected,
	}
}
