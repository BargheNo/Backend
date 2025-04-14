package enum

type TicketStatus uint

const (
	TicketStatusNotAnswered TicketStatus = iota + 1
	TicketStatusAnswered
	TicketStatusResolved
)

func (ts TicketStatus) String() string {
	switch ts {
	case TicketStatusNotAnswered:
		return "not answered"
	case TicketStatusAnswered:
		return "answered"
	case TicketStatusResolved:
		return "resolved"
	}
	return "unknown"
}
func GetAllTicketStatuses() []TicketStatus {
	return []TicketStatus{
		TicketStatusNotAnswered,
		TicketStatusAnswered,
		TicketStatusResolved,
	}
}
