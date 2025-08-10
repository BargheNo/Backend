package sortby

type TicketSortBy uint

const (
	TicketSortByCreatedAt TicketSortBy = iota + 1
)

func (s TicketSortBy) Name() string {
	switch s {
	case TicketSortByCreatedAt:
		return "تاریخ ایجاد"
	}
	return "ناشناس"
}

func (s TicketSortBy) DBColumn() string {
	switch s {
	case TicketSortByCreatedAt:
		return "created_at"
	}
	return ""
}

func GetTicketSortableColumns() map[TicketSortBy]bool {
	return map[TicketSortBy]bool{
		TicketSortByCreatedAt: true,
	}
}
