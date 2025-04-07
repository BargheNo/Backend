package enum

type UrgencyLevel uint

const (
	Low UrgencyLevel = iota + 1
	Medium
	High
)

func (s UrgencyLevel) String() string {
	switch s {
	case Low:
		return "low"
	case Medium:
		return "medium"
	case High:
		return "high"
	}
	return "unknown"
}
func GetAllUrgencyLevels() []UrgencyLevel {
	return []UrgencyLevel{
		Low,
		Medium,
		High,
	}
}
