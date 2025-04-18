package enum

type TicketSubject uint

const (
	TicketSubjectGeneral TicketSubject = iota + 1
	TicketSubjectPanel
	TicketSubjectInstallation
	TicketSubjectMaintenance
	TicketSubjectOther
)

func (s TicketSubject) String() string {
	switch s {
	case TicketSubjectGeneral:
		return "general"
	case TicketSubjectPanel:
		return "panel"
	case TicketSubjectInstallation:
		return "installation"
	case TicketSubjectMaintenance:
		return "maintenance"
	case TicketSubjectOther:
		return "other"
	}
	return "unknown"
}
func GetAllTicketSubjects() []TicketSubject {
	return []TicketSubject{
		TicketSubjectGeneral,
		TicketSubjectPanel,
		TicketSubjectInstallation,
		TicketSubjectMaintenance,
		TicketSubjectOther,
	}
}
