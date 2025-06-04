package enum

type UrgencyLevel uint

const (
	Low UrgencyLevel = iota + 1
	Medium
	High
)

func (level UrgencyLevel) String() string {
	switch level {
	case Low:
		return "عادی"
	case Medium:
		return "مهم"
	case High:
		return "فوری"
	}
	return "نامشخص"
}

func GetAllUrgencyLevels() []UrgencyLevel {
	return []UrgencyLevel{
		Low,
		Medium,
		High,
	}
}
