package enum

type StaffType uint

const (
	StaffTypeManager StaffType = iota + 1
	StaffTypeTechnician
	StaffTypeSupport
)

func (staff StaffType) String() string {
	switch staff {
	case StaffTypeManager:
		return "staffTypeManager"
	case StaffTypeTechnician:
		return "staffTypeTechnician"
	case StaffTypeSupport:
		return "staffTypeSupport"
	}
	return ""
}

func GetAllStaffTypes() []StaffType {
	return []StaffType{
		StaffTypeManager,
		StaffTypeTechnician,
		StaffTypeSupport,
	}
}
