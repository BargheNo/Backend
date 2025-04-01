package enum

type BidStatus uint

const (
	BidStatusPending BidStatus = iota + 1
	BidStatusAccepted
	BidStatusExpired
)

func (s BidStatus) String() string {
	switch s {
	case BidStatusPending:
		return "pending"
	case BidStatusAccepted:
		return "accepted"
	case BidStatusExpired:
		return "expired"
	}
	return "unknown"
}

func GetAllBidStatuses() []BidStatus {
	return []BidStatus{
		BidStatusPending,
		BidStatusAccepted,
		BidStatusExpired,
	}
}
