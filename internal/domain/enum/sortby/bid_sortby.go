package sortby

type BidSortBy uint

const (
	BidSortByCreatedAt BidSortBy = iota + 1
	BidSortByCost
	BidSortByInstallationTime
)

func (s BidSortBy) Name() string {
	switch s {
	case BidSortByCreatedAt:
		return "تاریخ ایجاد"
	case BidSortByCost:
		return "هزینه"
	case BidSortByInstallationTime:
		return "زمان نصب"
	}
	return "ناشناس"
}

func (s BidSortBy) DBColumn() string {
	switch s {
	case BidSortByCreatedAt:
		return "created_at"
	case BidSortByCost:
		return "cost"
	case BidSortByInstallationTime:
		return "installation_time"
	}
	return ""
}

func GetBidSortableColumns() map[BidSortBy]bool {
	return map[BidSortBy]bool{
		BidSortByCreatedAt:        true,
		BidSortByCost:             true,
		BidSortByInstallationTime: true,
	}
}
