package enum

type NewsStatus uint

const (
	NewsStatusActive NewsStatus = iota + 1
	NewsStatusDraft
)

func (status NewsStatus) String() string {
	switch status {
	case NewsStatusActive:
		return "active"
	case NewsStatusDraft:
		return "draft"
	}
	return ""
}

func GetAllNewsStatus() []NewsStatus {
	return []NewsStatus{
		NewsStatusActive,
		NewsStatusDraft,
	}
}
