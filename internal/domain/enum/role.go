package enum

type RoleName uint

const (
	SuperAdmin RoleName = iota + 1
	Customer
	Technician
	CorporationManager
	SupportAgent
	ContentManager
	Moderator
)

var rolePermissions = map[RoleName][]PermissionType{
	SuperAdmin: {
		PermissionAll,
	},
	Customer: {},
	Technician: {
		PanelViewUsageData, RepairViewAssigned, RepairAcceptRequest, RepairMarkComplete, ChatSendMessage, ChatViewAll,
	},
	CorporationManager: {
		CorporationManage, CorporationViewStats, CorporationUpdateSettings, PanelViewAssigned, PanelAssignToCustomer, PanelViewUsageData, PanelRemove, ReportViewOwn, AnalyticsViewAll,
		TicketCreate, TicketViewOwn, ReportViewOwn, ReportCreate, ProfileUpdate, UserViewCustomers,
	},
	SupportAgent: {
		TicketViewAll, TicketRespond, ReportViewAll, ReportManage, TicketCreate, ChatViewAll, UserViewAll, CorporationViewAll,
	},
	ContentManager: {
		BlogCreate, BlogEdit, BlogDelete, BlogView, BlogComment, NewsCreate, NewsEdit, NewsDelete, NewsView, NewsComment,
	},
	Moderator: {
		UserViewAll, TicketViewAll, BlogComment, NewsComment,
	},
}

func (role RoleName) Permissions() []PermissionType {
	if permissions, ok := rolePermissions[role]; ok {
		return permissions
	}
	return nil
}

func (role RoleName) String() string {
	switch role {
	case SuperAdmin:
		return "superAdmin"
	case Customer:
		return "customer"
	case Technician:
		return "technician"
	case CorporationManager:
		return "corporationManager"
	case SupportAgent:
		return "supportAgent"
	case ContentManager:
		return "contentManager"
	case Moderator:
		return "moderator"
	}
	return "unknown"
}

func GetAllRoleNames() []RoleName {
	return []RoleName{
		SuperAdmin,
		Customer,
		Technician,
		CorporationManager,
		SupportAgent,
		ContentManager,
		Moderator,
	}
}
