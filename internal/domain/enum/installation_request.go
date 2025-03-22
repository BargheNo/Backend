package enum

type InstallationRequestStatus uint

const (
	Active InstallationRequestStatus = iota + 1
	Expired
	Cancelled
)

func (status InstallationRequestStatus) String() string {
	switch status {
	case Active:
		return "active"
	case Expired:
		return "expired"
	case Cancelled:
		return "cancelled"
	}
	return ""
}

func GetAllBucketTypes() []InstallationRequestStatus {
	return []InstallationRequestStatus{
		Active,
		Expired,
		Cancelled,
	}
}
