package enum

type RequestType uint

const (
	InstallationRequest RequestType = iota + 1
	MaintenanceRequest
)

func (s RequestType) String() string {
	switch s {
	case InstallationRequest:
		return "installation_request"
	case MaintenanceRequest:
		return "maintenance_request"
	}
	return "unknown"
}

func GetAllRequestTypes() []RequestType {
	return []RequestType{
		InstallationRequest,
		MaintenanceRequest,
	}
}
