package enum

type TicketStatus uint

const (
	TicketStatusNotAnswered TicketStatus = iota + 1
	TicketStatusAnswered
	TicketStatusResolved
	TicketStatusAll
)

func (ts TicketStatus) String() string {
	switch ts {
	case TicketStatusNotAnswered:
		return "در انتظار پاسخ"
	case TicketStatusAnswered:
		return "پاسخ داده شده"
	case TicketStatusResolved:
		return "بسته شده"
	case TicketStatusAll:
		return "همه"
	}
	return "unknown"
}
func GetAllTicketStatuses() []TicketStatus {
	return []TicketStatus{
		TicketStatusNotAnswered,
		TicketStatusAnswered,
		TicketStatusResolved,
		TicketStatusAll,
	}
}
