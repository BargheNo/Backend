package enum

type GuaranteeStatus uint

const (
	GuaranteeStatusActive GuaranteeStatus = iota + 1
	GuaranteeStatusArchive
	GuaranteeStatusAll
)

func (status GuaranteeStatus) String() string {
	switch status {
	case GuaranteeStatusActive:
		return "فعال"
	case GuaranteeStatusArchive:
		return "آرشیو"
	case GuaranteeStatusAll:
		return "همه"
	}
	return "unknown"
}

func (status GuaranteeStatus) IsValid() bool {
	for _, validStatus := range GetAllGuaranteeStatuses() {
		if status == validStatus {
			return true
		}
	}
	return false
}

func GetAllGuaranteeStatuses() []GuaranteeStatus {
	return []GuaranteeStatus{
		GuaranteeStatusActive,
		GuaranteeStatusArchive,
		GuaranteeStatusAll,
	}
}
