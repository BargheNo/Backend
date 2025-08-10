package sortby

type InstallationSortBy uint

const (
	InstallationSortByCreatedAt InstallationSortBy = iota + 1
	InstallationSortByName
	InstallationSortByPowerRequest
	InstallationSortByMaxCost
)

func (s InstallationSortBy) Name() string {
	switch s {
	case InstallationSortByCreatedAt:
		return "تاریخ ایجاد"
	case InstallationSortByName:
		return "نام"
	case InstallationSortByPowerRequest:
		return "توان درخواستی"
	case InstallationSortByMaxCost:
		return "حداکثر هزینه"
	}
	return "ناشناس"
}

func (s InstallationSortBy) DBColumn() string {
	switch s {
	case InstallationSortByCreatedAt:
		return "created_at"
	case InstallationSortByName:
		return "name"
	case InstallationSortByPowerRequest:
		return "power_request"
	case InstallationSortByMaxCost:
		return "max_cost"
	}
	return ""
}

func GetInstallationSortableColumns() map[InstallationSortBy]bool {
	return map[InstallationSortBy]bool{
		InstallationSortByCreatedAt:    true,
		InstallationSortByName:         true,
		InstallationSortByPowerRequest: true,
		InstallationSortByMaxCost:      true,
	}
}
