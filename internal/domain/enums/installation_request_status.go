package enums

type InstallationRequestStatus uint

const (
	Open	InstallationRequestStatus = iota + 1
	Closed
)

func (s InstallationRequestStatus) String() string {
	switch s {
	case Open:
		return "open"
	case Closed:
		return "closed"
	}
	return "unknown"
}

func GetAllInstallationRequestStatuses() []InstallationRequestStatus {
	return []InstallationRequestStatus{
		Open,
		Closed,
	}
}