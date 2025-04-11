package enum

type MaintenanceRequestStatus uint

const (
	MaintenanceRequestStatusPending MaintenanceRequestStatus = iota + 1
	MaintenanceRequestStatusAccepted
	MaintenanceRequestStatusRejected
	MaintenanceRequestStatusCompleted
)

func (s MaintenanceRequestStatus) String() string {
	switch s {
	case MaintenanceRequestStatusPending:
		return "pending"
	case MaintenanceRequestStatusAccepted:
		return "accepted"
	case MaintenanceRequestStatusRejected:
		return "rejected"
	case MaintenanceRequestStatusCompleted:
		return "completed"
	}
	return "unknown"
}
func GetAllMaintenanceRequestStatuses() []MaintenanceRequestStatus {
	return []MaintenanceRequestStatus{
		MaintenanceRequestStatusPending,
		MaintenanceRequestStatusAccepted,
		MaintenanceRequestStatusRejected,
		MaintenanceRequestStatusCompleted,
	}
}
