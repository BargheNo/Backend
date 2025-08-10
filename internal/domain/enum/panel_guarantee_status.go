package enum

type PanelGuaranteeStatus uint

const (
	PanelGuaranteeStatusActive PanelGuaranteeStatus = iota + 1
	PanelGuaranteeStatusExpired
	PanelGuaranteeStatusVoided
	PanelGuaranteeStatusPending
	PanelGuaranteeStatusEmpty
)

func (status PanelGuaranteeStatus) String() string {
	switch status {
	case PanelGuaranteeStatusActive:
		return "فعال"
	case PanelGuaranteeStatusExpired:
		return "منقضی شده"
	case PanelGuaranteeStatusVoided:
		return "باطل شده"
	case PanelGuaranteeStatusPending:
		return "درانتظار نصب"
	case PanelGuaranteeStatusEmpty:
		return "بدون گارانتی"
	}
	return "unknown"
}

func GetAllPanelGuaranteeStatuses() []PanelGuaranteeStatus {
	return []PanelGuaranteeStatus{
		PanelGuaranteeStatusActive,
		PanelGuaranteeStatusExpired,
		PanelGuaranteeStatusVoided,
		PanelGuaranteeStatusPending,
		PanelGuaranteeStatusEmpty,
	}
}
