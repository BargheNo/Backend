package enum

type PanelStatus uint

const (
	PanelStatusActive = iota + 1
	PanelStatusPending
	PanelStatusDamaged
	PanelStatusAll
)

func (s PanelStatus) String() string {
	switch s {
	case PanelStatusActive:
		return "فعال"
	case PanelStatusPending:
		return "در انتظار نصب"
	case PanelStatusDamaged:
		return "خراب"
	case PanelStatusAll:
		return "همه"
	}
	return "unknown"
}

func GetAllPanelStatuses() []PanelStatus {
	return []PanelStatus{
		PanelStatusActive,
		PanelStatusPending,
		PanelStatusDamaged,
		PanelStatusAll,
	}
}
