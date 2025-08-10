package enum

type MaintenanceRequestStatus uint

const (
	MaintenanceRequestStatusPending MaintenanceRequestStatus = iota + 1
	MaintenanceRequestStatusAccepted
	MaintenanceRequestStatusRejected
	MaintenanceRequestStatusCompleted
	MaintenanceRequestStatusExpired
	MaintenanceRequestStatusCanceled
	MaintenanceRequestStatusAll
)

func (s MaintenanceRequestStatus) String() string {
	switch s {
	case MaintenanceRequestStatusPending:
		return "در انتظار تایید"
	case MaintenanceRequestStatusAccepted:
		return "تایید شده"
	case MaintenanceRequestStatusRejected:
		return "رد شده"
	case MaintenanceRequestStatusCompleted:
		return "تمام شده"
	case MaintenanceRequestStatusExpired:
		return "منقضی"
	case MaintenanceRequestStatusCanceled:
		return "لغو شده"
	case MaintenanceRequestStatusAll:
		return "همه"
	}
	return "unknown"
}

func GetAllowedMaintenanceRequestStatuses(role AgentType) []MaintenanceRequestStatus {
	switch role {
	case AgentTypeAdmin, AgentTypeCustomer:
		return []MaintenanceRequestStatus{
			MaintenanceRequestStatusPending,
			MaintenanceRequestStatusAccepted,
			MaintenanceRequestStatusRejected,
			MaintenanceRequestStatusCompleted,
			MaintenanceRequestStatusExpired,
			MaintenanceRequestStatusCanceled,
			MaintenanceRequestStatusAll,
		}
	case AgentTypeCorporation:
		return []MaintenanceRequestStatus{
			MaintenanceRequestStatusPending,
			MaintenanceRequestStatusAccepted,
			MaintenanceRequestStatusRejected,
			MaintenanceRequestStatusCompleted,
			MaintenanceRequestStatusAll,
		}
	default:
		return []MaintenanceRequestStatus{}
	}
}
