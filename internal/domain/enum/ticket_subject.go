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
		return "عمومی"
	case TicketSubjectPanel:
		return "پنل"
	case TicketSubjectInstallation:
		return "درخواست نصب"
	case TicketSubjectMaintenance:
		return "تعمیرات"
	case TicketSubjectOther:
		return "سایر"
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
