package enum

type NewsStatus uint

const (
	NewsStatusActive NewsStatus = iota + 1
	NewsStatusDraft
	NewsStatusAll
)

func (status NewsStatus) String() string {
	switch status {
	case NewsStatusActive:
		return "منتشر شده"
	case NewsStatusDraft:
		return "پیش نویس"
	case NewsStatusAll:
		return "همه"
	}
	return ""
}

func GetAllNewsStatus() []NewsStatus {
	return []NewsStatus{
		NewsStatusActive,
		NewsStatusDraft,
		NewsStatusAll,
	}
}
