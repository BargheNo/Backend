package enum

type InstallationRequestStatus uint

const (
	InstallationRequestStatusActive InstallationRequestStatus = iota + 1
	InstallationRequestStatusExpired
	InstallationRequestStatusCancelled
)

func (status InstallationRequestStatus) String() string {
	switch status {
	case InstallationRequestStatusActive:
		return "active"
	case InstallationRequestStatusExpired:
		return "expired"
	case InstallationRequestStatusCancelled:
		return "cancelled"
	}
	return ""
}

func GetAllInstallationRequests() []InstallationRequestStatus {
	return []InstallationRequestStatus{
		InstallationRequestStatusActive,
		InstallationRequestStatusExpired,
		InstallationRequestStatusCancelled,
	}
}
