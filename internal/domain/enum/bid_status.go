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
		return "در انتظار تایید"
	case BidStatusAccepted:
		return "تایید شده"
	case BidStatusExpired:
		return "منقضی"
	case BidStatusRejected:
		return "رد شده"
	case BidStatusCanceled:
		return "لغو شده"
	case BidStatusAll:
		return "همه"
	}
	return "unknown"
}

func GetAllBidStatuses() []BidStatus {
	return []BidStatus{
		BidStatusPending,
		BidStatusAccepted,
		BidStatusExpired,
		BidStatusRejected,
		BidStatusCanceled,
		BidStatusAll,
	}
}
