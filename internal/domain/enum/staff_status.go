package enum

type StaffStatus uint

const (
	StaffStatusActive StaffStatus = iota + 1
	StaffStatusInactive
	StaffStatusAll
)

func (s StaffStatus) String() string {
	switch s {
	case StaffStatusActive:
		return "فعال"
	case StaffStatusInactive:
		return "غیرفعال"
	case StaffStatusAll:
		return "همه"
	}
	return "نامعلوم"
}

func GetAllStaffStatuses() []StaffStatus {
	return []StaffStatus{
		StaffStatusActive,
		StaffStatusInactive,
		StaffStatusAll,
	}
}
