package enum

type InstallationRequestStatus uint

const (
	InstallationRequestStatusActive InstallationRequestStatus = iota + 1
	InstallationRequestStatusExpired
	InstallationRequestStatusCancelled
	InstallationRequestStatusDone
	InstallationRequestStatusAll
)

func (status InstallationRequestStatus) String() string {
	switch status {
	case InstallationRequestStatusActive:
		return "فعال"
	case InstallationRequestStatusExpired:
		return "منقضی"
	case InstallationRequestStatusCancelled:
		return "لغو شده"
	case InstallationRequestStatusDone:
		return "سپرده شده"
	case InstallationRequestStatusAll:
		return "همه"
	}
	return ""
}

func GetAllInstallationRequestStatuses() []InstallationRequestStatus {
	return []InstallationRequestStatus{
		InstallationRequestStatusActive,
		InstallationRequestStatusExpired,
		InstallationRequestStatusCancelled,
		InstallationRequestStatusDone,
		InstallationRequestStatusAll,
	}
}
