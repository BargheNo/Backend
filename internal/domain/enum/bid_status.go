package enum

type BidStatus uint

const (
	BidStatusPending BidStatus = iota + 1
	BidStatusAccepted
	BidStatusExpired
	BidStatusRejected
	BidStatusCanceled
	BidStatusAll
	BidStatusAllCustomer
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
	case BidStatusAllCustomer:
		return "همه"
	}
	return "unknown"
}

func GetUserBidStatuses() []BidStatus {
	return []BidStatus{
		BidStatusPending,
		BidStatusAccepted,
		BidStatusRejected,
		BidStatusAllCustomer,
	}
}

func GetCorporationBidStatuses() []BidStatus {
	return []BidStatus{
		BidStatusPending,
		BidStatusAccepted,
		BidStatusExpired,
		BidStatusRejected,
		BidStatusCanceled,
		BidStatusAll,
	}
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
