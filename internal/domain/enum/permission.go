package enum

type PermissionType uint

const (
	AccessAll PermissionType = iota + 1
	AccessCorporation
)

func (p PermissionType) String() string {
	switch p {
	case AccessAll:
		return "all"
	case AccessCorporation:
		return "corporation"
	}
	return ""
}

var permissionDescriptions = map[PermissionType]string{
	AccessAll:         "دسترسی کامل به سیستم",
	AccessCorporation: "دسترسی به مدیریت و اطلاعات شرکت",
}

func (p PermissionType) Description() string {
	if description, ok := permissionDescriptions[p]; ok {
		return description
	}
	return ""
}

func GetAllPermissionTypes() []PermissionType {
	return []PermissionType{
		AccessAll,
		AccessCorporation,
	}
}
